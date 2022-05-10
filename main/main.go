package main

import (
	"log"
	"net/http"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/route"
	"github.com/gorilla/handlers"
)

//go:generate go run main.go

func main() {
	router := route.GetRouter()
	methods := handlers.AllowedMethods([]string{http.MethodGet})
	log.Fatal(http.ListenAndServe(config.Port, handlers.CORS(methods)(router)))
}
