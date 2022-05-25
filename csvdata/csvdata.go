package csvdata

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
)

// func createCSVFile() {
// 	// Create the file
// 	dir, _ := os.Getwd()
// 	CSVfile, err := os.Create(CSVPath)
// 	if err != nil {
// 		log.Println("Err creating file:: ", err)
// 	}
// 	defer CSVfile.Close()
// }

// func copyResponseToCSVFile(resp *http.Response) {
// 	// Writer the body to file
// 	_, err = io.Copy(CSVfile, resp.Body)
// 	if err != nil {
// 		log.Println("Error copying body:: ", err)
// 	}
// 	csvFile, err := os.Open(CSVPath)
// 	if err != nil {
// 		log.Println("Error opening csv file:", err)
// 	}
// 	fmt.Println("Successfully Opened CSV file")
// 	defer csvFile.Close()
// }

func ExtractRowsFromCSVFile() (rows [][]string) {
	fmt.Println("Extracting linesssss....")
	csvFile, err := os.Open(config.CryptoNamesList)
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
	fmt.Println("READY TO PRINT ALL THE CRYPTOSSSS::::")
	for _, record := range csvReader {
		fmt.Println(record)
	}
	return csvReader
}
