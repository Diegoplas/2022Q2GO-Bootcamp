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
	"github.com/Diegoplas/2022Q2GO-Bootcamp/graph"

	"github.com/gorilla/mux"
)

// GraphCryptoRecords - Gets the historic data from http request, save it into a CSV file and graph of it.
func GraphCryptoRecords(w http.ResponseWriter, r *http.Request) {
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
	cryptoCodesRows, err := csvdata.ExtractRowsFromCSVFile(config.CryptoNamesList)
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
	err = csvdata.CopyResponseToCSVFile(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}

	// Use historical values data
	extractedHistoricalValuesRows, err := csvdata.ExtractRowsFromCSVFile(config.HistoricalValuesCSVPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	historicalValues, minPrice, maxPrice, err := csvdata.GetDataFromHistoricalValueRows(inputDays, extractedHistoricalValuesRows)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(string(err.Error())))
		return
	}
	graphName := graph.MakeGraph(historicalValues, minPrice, maxPrice, cryptoCode, rawInputDays)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Data successfully graphed in file: %s", graphName)))
}

// cryptoHystoricalValuesRequest - make http request.
func cryptoHystoricalValuesRequest(reqUrl string) (resp *http.Response, err error) {
	resp, err = http.Get(reqUrl)
	if err != nil {
		log.Println("HTTP request error: ", err)
		return nil, errors.New("request error")
	}
	return resp, nil
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
