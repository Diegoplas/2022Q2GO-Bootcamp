package usecase

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/model"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

type DataGetter interface {
	CreateCSVFile() error
	CopyResponseToCSVFile(resp *http.Response) error
	GetDataFromHistoricalValueRows(requestedDays int, historicalValueRows [][]string) (records model.CryptoRecordValues, dataError error)
	ExtractRowsFromCSVFile(csvFileName string) (rows [][]string, err error)
	ExtractDataFromBTCCSVRows(requestedDay int, csvRows [][]string) (records model.CryptoRecordValues, dataError error)
}

type DataCoverter interface {
	AverageHighLowCryptoPrices(lowPrice, highPrice string) (float64, error)
	ConvertCSVStrToDate(strDate string) (time.Time, error)
	ConvertCSVStrDataToNumericTypes(idStr, priceStr string) (int, float64, error)
}

type GraphMaker interface {
	MakeGraph(records model.CryptoRecordValues, cryptoCode, days string) string
}

type WorkPool struct {
	Getter    DataGetter
	Grapher   GraphMaker
	Converter DataCoverter
}

func NewWorkerPooler(Getter DataGetter, Grapher GraphMaker, Converter DataCoverter) WorkPool {
	return WorkPool{
		Getter:    Getter,
		Grapher:   Grapher,
		Converter: Converter,
	}
}

// GraphCryptoRecords - Gets the historic data from http request, save it into a CSV file and graph of it.
func (wp WorkPool) GraphCryptoRecords(requestedCryptoCode, requestedDays string) (cryptoGraphURL string, err error) {
	// Load config variables and Validate Input Crypto Code
	cryptoCodesRows, err := wp.Getter.ExtractRowsFromCSVFile(config.CryptoNamesListPath)
	if err != nil {
		return "", err
	}
	cryptoCode, err := validateInputCryptoCode(requestedCryptoCode, cryptoCodesRows)
	if err != nil {
		return "", err
	}
	requestURL, err := config.MakeRequestURL(requestedCryptoCode)
	if err != nil {
		return "", err
	}
	// Validate Input Num of Days to retrieve
	inputDays, err := validateInputDays(requestedDays)
	if err != nil {
		return "", err
	}
	// Get the data
	response, err := cryptoHystoricalValuesRequest(requestURL)
	if err != nil {
		return "", err
	}
	err = wp.Getter.CopyResponseToCSVFile(response)
	if err != nil {
		return "", err
	}
	// Use historical values data
	extractedHistoricalValuesRows, err := wp.Getter.ExtractRowsFromCSVFile(config.CryptoHistoricalValuesCSVPath)
	if err != nil {
		return "", err
	}
	historicalValues, err := wp.Getter.GetDataFromHistoricalValueRows(inputDays, extractedHistoricalValuesRows)
	if err != nil {
		return "", err
	}
	graphName := wp.Grapher.MakeGraph(historicalValues, cryptoCode, requestedDays)
	return graphName, nil
}

// cryptoHystoricalValuesRequest - make http request.
func cryptoHystoricalValuesRequest(reqUrl string) (resp *http.Response, err error) {
	// Get data from request
	request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		return nil, errors.New("request error")
	}
	response, err := Client.Do(request)
	if err != nil {
		if err != nil {
			log.Printf("error on the external api request: %v\n", err.Error())
			return nil, errors.New("response error")
		}
	}
	return response, nil
}

// GraphBTCValues - Gets the historic data from the CSV file and graph of it.
func (wp WorkPool) GraphBTCValues(inputDayStr string) (graphURL string, err error) {
	bitcoinCode := "BTC"
	inputDay, err := validateInputDays(inputDayStr)
	if err != nil {
		return "", err
	}
	extractedBTCHistoricalValuesRows, err := wp.Getter.ExtractRowsFromCSVFile(config.BTCHistoricalValuesCSVPath)
	if err != nil {
		return "", err
	}

	BTCRecordValues, err := wp.Getter.ExtractDataFromBTCCSVRows(inputDay, extractedBTCHistoricalValuesRows)
	if err != nil {
		return "", err
	}
	// If the requested day exceeds the number of records on the CSV file, use the latest day.
	if inputDay > len(extractedBTCHistoricalValuesRows) {
		inputDayStr = strconv.Itoa(len(extractedBTCHistoricalValuesRows))
	}
	graphName := wp.Grapher.MakeGraph(BTCRecordValues, bitcoinCode, inputDayStr)
	return graphName, nil
}

//validateInputDays - Validates the input is a valid positiv number.
func validateInputDays(inputDays string) (int, error) {
	inputDay, err := strconv.Atoi(inputDays)
	if err != nil {
		log.Println("error converting input string to int: ", err)
		err = errors.New("please insert a valid number")
		return 0, err
	}
	if inputDay < 1 || inputDay > 1001 {
		err = errors.New("the number of days you want the information about must be more than zero and less than 1001")
		return 0, err
	}
	return inputDay, nil
}

//validateInputCryptoName - Validates the input crypto name is contained in the available cryptos.
func validateInputCryptoCode(cryptoCode string, codesRows [][]string) (string, error) {
	cryptoCode = strings.ToUpper(cryptoCode)
	lenRows := len(codesRows)
	for idx := 1; idx < lenRows; idx++ {
		if cryptoCode == codesRows[idx][0] {
			return cryptoCode, nil
		}
	}
	return "", errors.New("please use a valid Crypto Currency code of cryptoCurrencyList.csv")
}

func (wp WorkPool) Worker(readChan chan []string, writeChan chan model.CryptoPricesAndDates, itemsPerWorker int) {
	worksCount := 0
	for {
		row, ok := <-readChan
		if !ok {
			return
		}

		strDate, strHighPrice, strLowPrice := row[0], row[2], row[3]
		// format the received row.
		date, err := wp.Converter.ConvertCSVStrToDate(strDate)
		if err != nil {
			log.Println(err)
			break
		}
		averagePrice, err := wp.Converter.AverageHighLowCryptoPrices(strHighPrice, strLowPrice)
		if err != nil {
			log.Println(err)
			break
		}

		cryptoRecordsInfo := model.CryptoPricesAndDates{
			Date:  date,
			Price: averagePrice,
		}

		writeChan <- cryptoRecordsInfo
		worksCount += 1

		// if worker completes it's corresponding items, the worker must rest.
		if worksCount == itemsPerWorker {
			break
		}
	}
}

// CSVWorkerPoolRowExtractor - Extract the information from a certain number of odd or even CSV rows concurrently with a given number of workers.
func (wp WorkPool) CSVWorkerPoolRowExtractor(oddOrEven string, items, itemsPerWorker int) ([]model.CryptoPricesAndDates, error) {
	// Validating input parameters.
	if itemsPerWorker > items {
		return nil, errors.New("number of items should be bigger than number of items per worker")
	}

	// Get the number of workers needed. If division is not exact, add another worker.
	numOfWorkers := items / itemsPerWorker
	if itemsRemainder := items % itemsPerWorker; itemsRemainder != 0 {
		numOfWorkers += 1
	}

	// Extract rows from CSV File
	csvRows, err := wp.Getter.ExtractRowsFromCSVFile(config.CryptoHistoricalValuesCSVPath)
	if err != nil || itemsPerWorker <= 0 {
		return nil, errors.New("data error")
	}
	if itemsPerWorker > items {
		return nil, errors.New("items should be more than items per worker")
	}

	// Create channels
	inputCh := make(chan []string)
	outputCh := make(chan model.CryptoPricesAndDates, items)

	// Waitgroup for Synchronization
	var waitGroup sync.WaitGroup

	// Declare the workers and send them their taks
	for workerID := 1; workerID <= numOfWorkers; workerID++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			wp.Worker(inputCh, outputCh, itemsPerWorker)
		}()
	}

	// Need to retrieve twice the requested items, because only will be getting half of them selecting odd/even
	totalItems := items * 2

	go func() {
		for idx := 1; idx <= totalItems; idx++ {
			if oddOrEven == "even" && (idx%2 == 0) {
				inputCh <- csvRows[idx]
			} else if oddOrEven == "odd" && (idx%2 == 1) {
				inputCh <- csvRows[idx]
			}
		}
		close(inputCh)
	}()

	// Wait workers to finish tasks
	waitGroup.Wait()
	close(outputCh)

	response := []model.CryptoPricesAndDates{}

	for requestedDatesAndPrices := range outputCh {
		response = append(response, requestedDatesAndPrices)
	}
	return response, nil
}
