package config

import "github.com/spf13/viper"

const (
	Port                    = ":8080"
	Market                  = "USD"
	HistoricalValuesCSVPath = "csvdata/historical-prices.csv"
	CryptoNamesList         = "csvdata/cryptoCurrencyList.csv"
	GraphTopBottomSpace     = 5000.0
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
