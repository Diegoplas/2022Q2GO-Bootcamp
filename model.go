package main

type CryptoData struct {
	MetaData struct {
		//TwoDigitalCurrencyCode   string `json:"2. Digital Currency Code"`
		ThreeDigitalCurrencyName string `json:"3. Digital Currency Name"`
		//FourMarketCode           string `json:"4. Market Code"`
		FiveMarketName string `json:"5. Market Name"`
	} `json:"Meta Data"`
}
