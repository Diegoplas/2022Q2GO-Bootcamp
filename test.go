package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	APIKey := os.Getenv("KEY")
	CSV_URL := fmt.Sprintf("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY_EXTENDED&symbol=IBM&interval=15min&slice=year1month1&apikey=%s", APIKey)

	// Create the file
	dir, _ := os.Getwd()
	CSVPath := dir + "/historical-prices.csv"
	CSVfile, err := os.Create(CSVPath)
	if err != nil {
		log.Println("Err creating file:: ", err)
	}
	defer CSVfile.Close()

	// Get the data
	resp, err := http.Get(CSV_URL)
	if err != nil {
		log.Println("Error with HTTP REQUEST:: ", err)
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(CSVfile, resp.Body)
	if err != nil {
		log.Println("Error copying body:: ", err)
	}

	csvFile, err := os.Open("historical-prices.csv")
	if err != nil {
		log.Println("Error opening csv file:", err)
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()
	csvReader, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Println("Error reading CSV data: ", err)
	}
	fmt.Println(csvReader)
	fmt.Println("READY TO PRINT::::")
	for _, record := range csvReader {
		fmt.Println(record)
	}
}
