package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/csvdata"

	"github.com/gorilla/mux"
)

// GraphCryptoRecords - Gets the historic data from the CSV file and graph of it.
func GraphCryptoRecords(w http.ResponseWriter, r *http.Request) {
	// APIKey := os.Getenv("KEY") // try viper
	// if APIKey == "" {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	w.Write([]byte("Please introduce global API Key"))
	// 	return
	// }

	// Validate Input Crypto Code
	inputCryptoCode := mux.Vars(r)["cryptoName"]
	cryptoCodesRows := csvdata.ExtractRowsFromCSVFile(config.CryptoNamesList)
	cryptoCode, err := validateInputCryptoCode(inputCryptoCode, cryptoCodesRows)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}
	// Validate Input Num of Days to retrieve
	inputDays := mux.Vars(r)["days"]
	fmt.Println(inputDays)
	days, err := validateInputDays(inputDays)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(string(err.Error())))
		return
	}
	log.Println("DAYSK:: ", days)
	// Get data from request

	//cryptoCurrency := "BTC"
	log.Println(cryptoCode)
	requestURL := fmt.Sprintf("https://www.alphavantage.co/query?function=DIGITAL_CURRENCY_DAILY&symbol=%s&market=%s&apikey=%s",
		cryptoCode, config.Market, APIKey)

	// Get the data
	response, err := cryptoHystoricalValuesRequest(requestURL)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(string(err.Error())))
		return
	}
	fmt.Println(response.Body)
	err = csvdata.CopyResponseToCSVFile(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}

	// Use historical values data
	historicalValuesRows := csvdata.ExtractRowsFromCSVFile(config.HistoricalValuesCSVPath)
	//mandams graficar y los datos como en la primera entrega (agregar days)
	//// ADD VALIDATION FOR INPUT

	// BTCRecords := model.CryptoRecordValues{}
	// BTCRecords, minValue, maxValue, err := csvdata.ExtractFromCSV(inputDay)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// 	return
	// }
	// graph.MakeGraph(BTCRecords, minValue, maxValue)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Data successfully graphed in file: %s", config.PNGFileName)))
}

func cryptoHystoricalValuesRequest(reqUrl string) (resp *http.Response, err error) {
	// Get the data
	fmt.Println(reqUrl)
	resp, err = http.Get(reqUrl)
	if err != nil {
		log.Println("HTTP request error: ", err)
		return nil, errors.New("request error")
	}
	return resp, nil
}

//validateInputDays - Validates the input is a valid positiv number.
func validateInputDays(inputDays string) (int, error) {
	fmt.Println(inputDays)
	inputDay, err := strconv.Atoi(inputDays)
	if err != nil {
		log.Println("error converting input string to int: ", err)
		err = errors.New("please insert a valid number")
		return 0, err
	}
	if inputDay < 1 {
		err = errors.New("number must be more than zero")
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
