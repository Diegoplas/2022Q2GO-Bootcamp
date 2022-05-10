package controller

import (
	"net/http"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/csvdata"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/graph"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/model"
	"github.com/gorilla/mux"
)

func GraphBTCValues(w http.ResponseWriter, r *http.Request) {
	requestedDate := mux.Vars(r)["date"]
	BTCRecords := model.CryptoValueRecords{}
	BTCRecords, minValue, maxValue := csvdata.ExtractFromCSV(requestedDate)
	graph.MakeGraph(BTCRecords, minValue, maxValue)
}
