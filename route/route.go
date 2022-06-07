package route

import (
	"net/http"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/controller"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/csvdata"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/graph"

	"github.com/gorilla/mux"
)

func GetRouter() (router *mux.Router) {

	dataService := csvdata.NewCSVDataHandler()
	graphService := graph.NewGrapher()
	converterService := csvdata.NewCSVDataConverter()
	dataHandler := controller.NewDataGetter(dataService, graphService, converterService)

	router = mux.NewRouter()
	router.HandleFunc("/btc-values/{day}", dataHandler.GraphBTCValues).Methods(http.MethodGet)
	router.HandleFunc("/usd-crypto-conversion/{cryptoCode}/{days}", dataHandler.GraphCryptoRecords).Methods(http.MethodGet)
	router.HandleFunc("/workerpool/{odd_or_even}/{items}/{items_per_worker}", dataHandler.WorkerPoolHandler).Methods(http.MethodGet)
	return router
}
