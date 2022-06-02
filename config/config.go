package config

import "github.com/spf13/viper"

const (
	Port                          = ":8080"
	Market                        = "USD"
	CryptoNamesListPath           = "csvdata/cryptoCurrencyList.csv"
	CryptoHistoricalValuesCSVPath = "csvdata/crypto-historical-prices.csv"
	BTCHistoricalValuesCSVPath    = "csvdata/BTC-historical-prices.csv"
)

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
