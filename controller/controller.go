package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/csvdata"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/graph"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/model"

	"github.com/gorilla/mux"
)

// GraphBTCValues - Gets the historic data from the CSV file and graph of it.
func GraphBTCValues(w http.ResponseWriter, r *http.Request) {
	requestedDay := mux.Vars(r)["day"]
	fmt.Println(requestedDay)
	inputDay, err := validInputDay(requestedDay)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please insert a valid positive number."))
		return
	}
	BTCRecords := model.CryptoValueRecords{}
	BTCRecords, minValue, maxValue, err := csvdata.ExtractFromCSV(inputDay)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	graph.MakeGraph(BTCRecords, minValue, maxValue)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Data successfully graphed in file: %s", config.PNGFileName)))

}

// validInputDay - Validates the input is a valid positiv number.
func validInputDay(input string) (int, error) {
	fmt.Println(input)
	inputDay, err := strconv.Atoi(input)
	if err != nil {
		log.Println("error converting input string to int: ", err)
		return 0, err
	}
	if inputDay < 1 {
		err = errors.New("number must be more than zero")
		return 0, err
	}
	return inputDay, nil
}
