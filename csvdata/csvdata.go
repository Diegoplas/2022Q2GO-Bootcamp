package csvdata

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
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

func ExtractRowsFromCSVFile(csvFileName string) (rows [][]string) {
	csvFile, err := os.Open(csvFileName)
	if err != nil {
		log.Println("Error opening csv file:", err)
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	CSVRows, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Println("Error reading CSV data: ", err)
	}
	return CSVRows
}
