package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/model"

	"github.com/gorilla/mux"
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
	ExtractRowsFromCSVFile(csvFileName string) (rows [][]string, err error)
	ExtractDataFromBTCCSVRows(requestedDay int, csvRows [][]string) (records model.CryptoRecordValues, dataError error)
	GetDataFromHistoricalValueRows(requestedDays int, historicalValueRows [][]string) (records model.CryptoRecordValues, dataError error)
}

type GraphMaker interface {
	MakeGraph(records model.CryptoRecordValues, cryptoCode, days string) string
}

type DataHandlerAndGrapher struct {
	getter  DataGetter
	grapher GraphMaker
}

func NewDataGetter(getter DataGetter, grapher GraphMaker) DataHandlerAndGrapher {
	return DataHandlerAndGrapher{
		getter:  getter,
		grapher: grapher,
	}
}

// GraphCryptoRecords - Gets the historic data from http request, save it into a CSV file and graph of it.
func (dhg DataHandlerAndGrapher) GraphCryptoRecords(w http.ResponseWriter, r *http.Request) {
	// Load config variables
	configVars, err := config.LoadConfig(".")
	if err != nil {
		log.Println("cannot load config:", err)
	}
	if configVars.APIKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please introduce use a valid API Key"))
		return
	}
	// Validate Input Crypto Code
	rawInputCryptoCode := mux.Vars(r)["cryptoCode"]
	cryptoCodesRows, err := dhg.getter.ExtractRowsFromCSVFile(config.CryptoNamesListPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	cryptoCode, err := validateInputCryptoCode(rawInputCryptoCode, cryptoCodesRows)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}
	// Validate Input Num of Days to retrieve
	rawInputDays := mux.Vars(r)["days"]
	inputDays, err := validateInputDays(rawInputDays)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}
	// Get data from request
	requestURL := fmt.Sprintf("https://www.alphavantage.co/query?function=DIGITAL_CURRENCY_DAILY&symbol=%s&market=%s&apikey=%s&datatype=csv",
		cryptoCode, config.Market, configVars.APIKey)
	// Get the data
	response, err := cryptoHystoricalValuesRequest(requestURL)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(string(err.Error())))
		return
	}
	err = dhg.getter.CopyResponseToCSVFile(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	// Use historical values data
	extractedHistoricalValuesRows, err := dhg.getter.ExtractRowsFromCSVFile(config.CryptoHistoricalValuesCSVPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	historicalValues, err := dhg.getter.GetDataFromHistoricalValueRows(inputDays, extractedHistoricalValuesRows)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	graphName := dhg.grapher.MakeGraph(historicalValues, cryptoCode, rawInputDays)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Data successfully graphed in file: %s", graphName)))
}

// GraphBTCValues - Gets the historic data from the CSV file and graph of it.
func (dhg DataHandlerAndGrapher) GraphBTCValues(w http.ResponseWriter, r *http.Request) {
	BitcoinCode := "BTC"
	rawInputDay := mux.Vars(r)["day"]
	inputDay, err := validateInputDays(rawInputDay)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}
	extractedBTCHistoricalValuesRows, err := dhg.getter.ExtractRowsFromCSVFile(config.BTCHistoricalValuesCSVPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	BTCRecordValues, err := dhg.getter.ExtractDataFromBTCCSVRows(inputDay, extractedBTCHistoricalValuesRows)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	// If the requested day exceeds the number of records on the CSV file, use the latest day.
	if inputDay > len(extractedBTCHistoricalValuesRows) {
		rawInputDay = strconv.Itoa(len(extractedBTCHistoricalValuesRows))
	}
	graphName := dhg.grapher.MakeGraph(BTCRecordValues, BitcoinCode, rawInputDay)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Data successfully graphed in file: %s", graphName)))
}

// cryptoHystoricalValuesRequest - make http request.
func cryptoHystoricalValuesRequest(reqUrl string) (resp *http.Response, err error) {
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
