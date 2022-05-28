package route

import (
	"net/http"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/controller"

	"github.com/gorilla/mux"
)

func GetRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/usd-crypto-conversion/{cryptoCode}/{days}", controller.GraphCryptoRecords).Methods(http.MethodGet)
	return router
}
