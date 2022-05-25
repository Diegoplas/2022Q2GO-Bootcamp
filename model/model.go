package model

import "time"

type CryptoInformation struct {
	MetaData struct {
		//TwoDigitalCurrencyCode   string `json:"2. Digital Currency Code"`
		ThreeDigitalCurrencyName string `json:"3. Digital Currency Name"`
		//FourMarketCode           string `json:"4. Market Code"`
		FiveMarketName string `json:"5. Market Name"`
	} `json:"Meta Data"`
}

type CryptoRecordValues struct {
	Ids    []int       //intentar sacar como int
	Dates  []time.Time //intentar sacar como date
	Values []float64   //intentar sacar como float
}
