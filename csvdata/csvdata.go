package csvdata

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/model"

	"github.com/wcharczuk/go-chart"
)

type DataError struct {
}

func (de *DataError) Error() string {
	return "error handling data"
}

type CSVDataHandler struct {
}

func NewCSVDataHandler() CSVDataHandler {
	return CSVDataHandler{}
}

func (csvdh CSVDataHandler) CreateCSVFile() error {
	// Create the file
	_, err := os.Create(config.CryptoHistoricalValuesCSVPath)
	if err != nil {
		log.Println("Error creating csvFile ", err) //CHECK FOR ERROR
		return &DataError{}
	}
	return nil
}

// CopyResponseToCSVFile - copy the response body of a request to just created CSV file.
func (csvdh CSVDataHandler) CopyResponseToCSVFile(resp *http.Response) error {
	// Writer the body to file
	err := csvdh.CreateCSVFile()
	if err != nil {
		return err
	}
	CSVFile, err := os.OpenFile(config.CryptoHistoricalValuesCSVPath, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Println("Error opening csv file:", err)
		return &DataError{}
	}
	_, err = io.Copy(CSVFile, resp.Body)
	if err != nil {
		log.Println("Error copying body: ", err)
		return &DataError{}
	}
	defer CSVFile.Close()
	defer resp.Body.Close()
	return nil
}

// ExtractRowsFromCSVFile - Reads the rows of a CSV File and returns a 2D array (row and columns)
func (csvdh CSVDataHandler) ExtractRowsFromCSVFile(csvFileName string) (rows [][]string, err error) {
	csvFile, err := os.Open(csvFileName)
	if err != nil {
		log.Println("Error opening csv file:", err)
		return nil, &DataError{}
	}
	defer csvFile.Close()
	CSVRows, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Println("Error reading CSV data: ", err)
		return nil, &DataError{}
	}
	return CSVRows, nil
}

// GetDataFromHistoricalValueRows - Convert the data extracted from the CSV to their corresponding types, store them and get the Min and Max price values.
func (csvdh CSVDataHandler) GetDataFromHistoricalValueRows(requestedDays int, historicalValueRows [][]string) (records model.CryptoRecordValues, dataError error) {
	cryptoRecords := model.CryptoRecordValues{}
	// Default values for obtaining the max and the min values (for graph use).
	minValue, maxValue := 100000000.0, 0.0
	compensateRowTitles := requestedDays + 1
	// Starting from 1 to skip the column titles and adding one to the request day to compensate.
	for idx := 1; idx <= compensateRowTitles; idx++ {
		date := historicalValueRows[idx][0]
		highPriceUSD := historicalValueRows[idx][2]
		lowPriceUSD := historicalValueRows[idx][3]
		timestamp, err := convertCSVStrToDate(date)
		if err != nil {
			return model.CryptoRecordValues{}, errors.New("historical data error")
		}
		value, err := averageHighLowCryptoPrices(lowPriceUSD, highPriceUSD)
		if err != nil {
			return model.CryptoRecordValues{}, errors.New("historical data error")
		}
		if value > maxValue {
			maxValue = value
		}
		if value < minValue {
			minValue = value
		}
		cryptoRecords.Dates = append(cryptoRecords.Dates, timestamp)
		cryptoRecords.AveragePrice = append(cryptoRecords.AveragePrice, value)
	}
	cryptoRecords.MaxPrice = maxValue
	cryptoRecords.MinPrice = minValue
	return cryptoRecords, nil
}

func (csvdh CSVDataHandler) ExtractDataFromBTCCSVRows(requestedDay int, csvRows [][]string) (records model.CryptoRecordValues, dataError error) {
	BTCRecords := model.CryptoRecordValues{}
	// Default values for obtaining the max and the min values (for graph use).
	minValue, maxValue := 1000000.0, 0.0
	csvLinesLen := len(csvRows)
	// Starting from 1 to skip the column titles.
	for idx := 1; idx < csvLinesLen; idx++ {
		dateOnly := csvRows[idx][1][:10]
		date, err := convertCSVStrToDate(dateOnly)
		if err != nil {
			return model.CryptoRecordValues{}, err
		}
		id, value, err := convertCSVStrDataToNumericTypes(csvRows[idx][0], csvRows[idx][2])
		if err != nil {
			return model.CryptoRecordValues{}, err
		}
		if value > maxValue {
			maxValue = value
		}
		if value < minValue {
			minValue = value
		}
		BTCRecords.Ids = append(BTCRecords.Ids, id)
		BTCRecords.Dates = append(BTCRecords.Dates, date)
		BTCRecords.AveragePrice = append(BTCRecords.AveragePrice, value)
		if requestedDay == id {
			break
		}
	}
	BTCRecords.MaxPrice = maxValue
	BTCRecords.MinPrice = minValue
	return BTCRecords, nil
}

func convertCSVStrDataToNumericTypes(strId, strFloat string) (int, float64, error) {
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Println("Error, converting id (int):", err)
		return 0, 0.0, errors.New("data error, numeric")
	}
	value, err := strconv.ParseFloat(strFloat, 64)
	if err != nil {
		log.Println("Error, converting value (float):", err)
		return 0, 0.0, errors.New("data error, numeric")
	}
	return id, value, nil
}

// averageHighLowCryptoPrices - Returns the average of the low and high price of the crypto currency.
func averageHighLowCryptoPrices(lowPrice, highPrice string) (float64, error) {
	lowP, err := strconv.ParseFloat(lowPrice, 64)
	if err != nil {
		log.Println("Error, converting value (float):", err)
		return 0.0, errors.New("data error, numeric")
	}
	highP, err := strconv.ParseFloat(highPrice, 64)
	if err != nil {
		log.Println("Error, converting value (float):", err)
		return 0.0, errors.New("data error, numeric")
	}
	averagePrice := (lowP + highP) / 2
	return averagePrice, nil
}

// convertCSVStrToDate - Convert string date to time.Time type
func convertCSVStrToDate(strDate string) (time.Time, error) {
	date, err := time.Parse(chart.DefaultDateFormat, strDate)
	if err != nil {
		log.Println("Error, date parsing:", err)
		return time.Time{}, errors.New("data error, date")
	}
	return date, nil
}
