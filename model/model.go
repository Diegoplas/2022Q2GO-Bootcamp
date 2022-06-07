package model

import "time"

type CryptoRecordValues struct {
	Ids          []int
	Dates        []time.Time
	AveragePrice []float64
	MaxPrice     float64
	MinPrice     float64
}

type CryptoPricesAndDates struct {
	Date  time.Time
	Price float64
}
