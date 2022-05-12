package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Diegoplas/2022Q2GO-Bootcamp/config"
	"github.com/Diegoplas/2022Q2GO-Bootcamp/route"
	"github.com/gorilla/handlers"
)

//go:generate go run main.go
const (
	welcomeString = `:::::::: BTC-USD Data Grapher ::::::::
::::::::::::::::::::::::::::::::::::::
:     Please use as an input the     :
:    day since you want to obtain    :
:      the prices, eg 7 are the      :
: historical prices of the last week :
::::::::::::::::::::::::::::::::::::::`
)

func main() {
	fmt.Println(welcomeString)
	router := route.GetRouter()
	methods := handlers.AllowedMethods([]string{http.MethodGet})
	log.Fatal(http.ListenAndServe(config.Port, handlers.CORS(methods)(router)))
}
