package controller

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/csvdata"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/graph"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/model"
	"github.com/gorilla/mux"
)

func GraphBTCValues(w http.ResponseWriter, r *http.Request) {
	requestedDate := mux.Vars(r)["date"]
	if !validInputDate(requestedDate) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Date format must: YYYY-MM-DD (hyphens included)"))
		return
	}
	BTCRecords := model.CryptoValueRecords{}
	BTCRecords, minValue, maxValue, err := csvdata.ExtractFromCSV(requestedDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	graph.MakeGraph(BTCRecords, minValue, maxValue)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Data successfully graphed in file: %s", config.PNGFileName)))

}

// Format must be YYYY-MM-DD
func validInputDate(inputDate string) bool {
	defaultDateFormat := "2006-01-02"
	_, err := time.Parse(defaultDateFormat, inputDate)
	if err != nil {
		log.Println("Error, input date:", err)
		return false
	}
	return true
}
