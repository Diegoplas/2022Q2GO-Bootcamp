package route

import (
	"net/http"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/controller"
)

func GetRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/btc-values/{date}", controller.GraphBTCValues).Methods(http.MethodGet)
	return router
}
