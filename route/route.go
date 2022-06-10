package route

import (
	"net/http"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/controller"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/csvdata"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/graph"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/usecase"

	"github.com/gorilla/mux"
)

func GetRouter() (router *mux.Router) {

	dataService := csvdata.NewCSVDataHandler()
	graphService := graph.NewGrapher()
	converterService := csvdata.NewCSVDataConverter()
	workerPoolService := usecase.NewWorkerPooler(dataService, graphService, converterService)
	serviceHandler := controller.NewDataGrapher(workerPoolService)

	router = mux.NewRouter()
	router.HandleFunc("/btc-values/{day}", serviceHandler.GraphBTCValuesHandler).Methods(http.MethodGet)
	router.HandleFunc("/usd-crypto-conversion/{cryptoCode}/{days}", serviceHandler.GraphCryptoRecordsHandler).Methods(http.MethodGet)
	router.HandleFunc("/workerpool/{odd_or_even}/{items}/{items_per_worker}", serviceHandler.WorkerPoolHandler).Methods(http.MethodGet)
	return router
}
