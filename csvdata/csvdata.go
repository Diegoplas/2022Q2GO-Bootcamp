package csvdata

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/model"
	"github.com/wcharczuk/go-chart"
)

func readCSVData() ([][]string, error) {
	csvFile, err := os.Open(config.CSVFileName)
	if err != nil {
		log.Println("Error opening csv file:", err)
		return nil, errors.New("error reading data")
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Println("Error reading information:", err)
		return nil, errors.New("error reading data")
	}
	return csvLines, nil
}
func ExtractFromCSV(requestedDate string) (records model.CryptoValueRecords, minValue, maxValue float64, dataError error) {
	BTCRecords := model.CryptoValueRecords{}
	csvLines, err := readCSVData()
	if err != nil {
		return model.CryptoValueRecords{}, 0, 0, err
	}

	// Check if the input date is earlier than the latest date of the CSV file.
	latestDateCSVData := csvLines[1][1]
	if !validDate(latestDateCSVData, requestedDate) {
		return model.CryptoValueRecords{}, 0, 0, errors.New("please use an older date")
	}

	// Default values for obtaining the max and the min values (for graph use).
	minValue, maxValue = 1000000.0, 0.0
	csvLinesLen := len(csvLines)
	// Starting from 1 to skip the column titles.
	for idx := 1; idx < csvLinesLen; idx++ {
		dateOnly := csvLines[idx][1][:10]
		date, err := convertStrToDate(dateOnly)
		if err != nil {
			return model.CryptoValueRecords{}, 0.0, 0.0, err
		}
		id, value, err := convertStrDataToNumericTypes(csvLines[idx][0], csvLines[idx][2])
		if err != nil {
			return model.CryptoValueRecords{}, 0.0, 0.0, err
		}
		if value > maxValue {
			maxValue = value
		}
		if value < minValue {
			minValue = value
		}
		BTCRecords.Ids = append(BTCRecords.Ids, id)
		BTCRecords.Dates = append(BTCRecords.Dates, date)
		BTCRecords.Values = append(BTCRecords.Values, value)
		if requestedDate == dateOnly {
			break
		}
	}

	return BTCRecords, minValue, maxValue, nil
}

func validDate(latestDateCSVData, userRequestedDate string) bool {
	// Get only the date of the string
	latestDateCSVData = latestDateCSVData[:10]
	latestDate, err := convertStrToDate(latestDateCSVData)
	if err != nil {
		return false
	}
	requestedDate, err := convertStrToDate(userRequestedDate)
	if err != nil {
		return false
	}
	if requestedDate.After(latestDate) {
		return false
	}
	return true
}

func convertStrDataToNumericTypes(strId, strFloat string) (int, float64, error) {
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

func convertStrToDate(strDate string) (time.Time, error) {
	date, err := time.Parse(chart.DefaultDateFormat, strDate)
	if err != nil {
		log.Println("Error, date parsing:", err)
		return time.Time{}, errors.New("data error, date")
	}
	return date, nil
}
