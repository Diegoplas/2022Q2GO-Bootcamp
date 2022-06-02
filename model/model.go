package model

import "time"

type CryptoRecordValues struct {
	Ids          []int
	Dates        []time.Time
	AveragePrice []float64
	MaxPrice     float64
	MinPrice     float64
}
