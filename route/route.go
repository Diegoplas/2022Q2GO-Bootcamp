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
	dataHandler := controller.NewDataGetter(dataService, graphService)

	router = mux.NewRouter()
	router.HandleFunc("/usd-crypto-conversion/{cryptoCode}/{days}", dataHandler.GraphCryptoRecords).Methods(http.MethodGet)
	return router
}
