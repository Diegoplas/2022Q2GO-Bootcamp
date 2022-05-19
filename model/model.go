package model

import "time"

type CryptoValueRecords struct {
	Ids    []int       //intentar sacar como int
	Dates  []time.Time //intentar sacar como date
	Values []float64   //intentar sacar como float
}
