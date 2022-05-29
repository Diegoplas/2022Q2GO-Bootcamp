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

func createCSVFile() error {
	// Create the file
	_, err := os.Create(config.HistoricalValuesCSVPath)
	if err != nil {
		log.Println("Error creating csvFile ", err) //CHECK FOR ERROR
		return &DataError{}
	}
	return nil
}

// CopyResponseToCSVFile - copy the response body of a request to just created CSV file.
func CopyResponseToCSVFile(resp *http.Response) error {
	// Writer the body to file
	err := createCSVFile()
	if err != nil {
		return err
	}
	CSVFile, err := os.OpenFile(config.HistoricalValuesCSVPath, os.O_RDWR|os.O_APPEND, 0660)
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
func ExtractRowsFromCSVFile(csvFileName string) (rows [][]string, err error) {
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
func GetDataFromHistoricalValueRows(requestedDays int, historicalValueRows [][]string) (records model.CryptoRecordValues, minValue, maxValue float64, dataError error) {
	cryptoRecords := model.CryptoRecordValues{}
	// Default values for obtaining the max and the min values (for graph use).
	minValue, maxValue = 100000000.0, 0.0
	// Starting from 1 to skip the column titles and adding one to the request day to compensate.
	for idx := 1; idx <= requestedDays+1; idx++ {
		date := historicalValueRows[idx][0]
		highPriceUSD := historicalValueRows[idx][2]
		lowPriceUSD := historicalValueRows[idx][3]
		timestamp, err := convertCSVStrToDate(date)
		if err != nil {
			return model.CryptoRecordValues{}, 0.0, 0.0, errors.New("historical data error")
		}
		value, err := averageHighLowCryptoPrices(lowPriceUSD, highPriceUSD)
		if err != nil {
			return model.CryptoRecordValues{}, 0.0, 0.0, errors.New("historical data error")
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
	return cryptoRecords, minValue, maxValue, nil
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
