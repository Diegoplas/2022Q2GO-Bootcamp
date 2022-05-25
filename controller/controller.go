package controller

import (
	"fmt"
	"net/http"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/csvdata"

	"github.com/gorilla/mux"
)

// GraphCryptoRecords - Gets the historic data from the CSV file and graph of it.
func GraphCryptoRecords(w http.ResponseWriter, r *http.Request) {
	// APIKey := os.Getenv("KEY")
	// if APIKey == "" {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	w.Write([]byte("Please introduce global API Key"))
	// 	return
	// }
	name := mux.Vars(r)["cryptoName"]
	fmt.Println(name)
	csvdata.ExtractRowsFromCSVFile()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Data successfully graphed in file: %s", config.PNGFileName)))
	//// ADD VALIDATION FOR INPUT
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("Please insert a valid positive number."))
	// 	return
	// }
	// BTCRecords := model.CryptoRecordValues{}
	// BTCRecords, minValue, maxValue, err := csvdata.ExtractFromCSV(inputDay)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// graph.MakeGraph(BTCRecords, minValue, maxValue)

}

// func cryptoHystoricalValuesRequest() (resp *http.Response) {
// 	// Get the data
// 	resp, err := http.Get(config.RequestURL)
// 	if err != nil {
// 		log.Println("Error with HTTP REQUEST:: ", err)
// 	}
// 	defer resp.Body.Close()
// 	return resp
// }

// func validateRequestURL(cryptoCurrency, APIKey string) string {
// 	upperCrypto := strings.ToUpper(cryptoCurrency)
// 	requestURL := fmt.Sprintf("https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=%s&to_currency=CNY&apikey=demo%s",
// 		upperCrypto, APIKey)
// 	// añadir lista extraida del csv y ver si existe esa crypto
// 	return requestURL
// }

// validateInputDay - Validates the input is a valid positiv number.
// func validateInputDay(input string) (int, error) {
// 	fmt.Println(input)
// 	inputDay, err := strconv.Atoi(input)
// 	if err != nil {
// 		log.Println("error converting input string to int: ", err)
// 		return 0, err
// 	}
// 	if inputDay < 1 {
// 		err = errors.New("number must be more than zero")
// 		return 0, err
// 	}
// 	return inputDay, nil
// }
