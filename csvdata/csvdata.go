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

func ExtractFromCSV(requestedDate string) (records model.CryptoValueRecords, minValue, maxValue float64) {
	BTCRecords := model.CryptoValueRecords{}
	csvFile, err := os.Open(config.CSVFileName)
	if err != nil {
		log.Println("Error opening csv file:", err)
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Println("Error reading information:", err)
	}
	csvLinesLen := len(csvLines)
	//default values for obtaining the max and the min values.
	minValue, maxValue = 1000000.0, 0.0
	//Starting from 1 to skip the column titles.
	for idx := 1; idx < csvLinesLen; idx++ {
		dateOnly := csvLines[idx][1][:10]
		id, date, value, err := convertCSVDataToTypes(csvLines[idx][0], dateOnly, csvLines[idx][2])
		if err != nil {
			return model.CryptoValueRecords{}, 0.0, 0.0
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
	return BTCRecords, minValue, maxValue
}

// Convert CSV string data to their corresponding types.
func convertCSVDataToTypes(strId, strDate, strFloat string) (int, time.Time, float64, error) {
	id, err := strconv.Atoi(strId)
	if err != nil {
		log.Println("Error, converting id (int):", err)
		return 0, time.Time{}, 0.0, errors.New("data error")
	}

	timeInput := strDate
	date, err := time.Parse(chart.DefaultDateFormat, timeInput)
	if err != nil {
		log.Println("Error, date parsing:", err)
		return 0, time.Time{}, 0.0, errors.New("data error")
	}

	value, err := strconv.ParseFloat(strFloat, 64)
	if err != nil {
		log.Println("Error, converting value (float):", err)
		return 0, time.Time{}, 0.0, errors.New("data error")
	}
	return id, date, value, err
}
