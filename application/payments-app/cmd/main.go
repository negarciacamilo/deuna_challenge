package main

import (
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/http"
	"github.com/spf13/viper"
	"os"
)

func main() {
	if env := os.Getenv("ENVIRONMENT"); env == "" {
		viper.SetConfigFile("env.json")
	} else {
		viper.SetConfigFile("dockerenv.json")
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err = http.NewRouter().Run(viper.GetString("PAYMENTS_PORT")); err != nil {
		panic(err)
	}
}
