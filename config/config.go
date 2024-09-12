package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const (
	Port                          = ":8080"
	CryptoNamesListPath           = "csvdata/cryptoCurrencyList.csv"
	CryptoHistoricalValuesCSVPath = "csvdata/crypto-historical-prices.csv"
	BTCHistoricalValuesCSVPath    = "csvdata/BTC-historical-prices.csv"
)

func MakeRequestURL(cryptoCode string) (url string, err error) {
	configVars, err := LoadConfig(".")
	if err != nil {
		log.Println("cannot load config:", err)
	}
	if configVars.APIKey == "" {
		return "", errors.New("please introduce use a valid API Key")
	}

	requestURL := fmt.Sprintf("https://www.alphavantage.co/query?function=DIGITAL_CURRENCY_DAILY&symbol=%s&market=USD&apikey=%s&datatype=csv",
		cryptoCode, configVars.APIKey)

	return requestURL, nil
}

type environmentVars struct {
	APIKey string `mapstructure:"API_KEY"`
}

func LoadConfig(configPath string) (envVars environmentVars, err error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&envVars)
	return
}
