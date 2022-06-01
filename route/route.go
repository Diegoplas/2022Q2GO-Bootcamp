package route

import (
	"net/http"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/controller"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/csvdata"

	"github.com/gorilla/mux"
)

func GetRouter() (router *mux.Router) {

	//dataService := csvdata.NewRepositoryService(csvdata.PokemonRepo{})
	dataService := csvdata.NewCSVDataHandler()
	dataHandler := controller.NewDataGetter(dataService)

	router = mux.NewRouter()
	router.HandleFunc("/usd-crypto-conversion/{cryptoCode}/{days}", dataHandler.GraphCryptoRecords).Methods(http.MethodGet)
	return router
}
