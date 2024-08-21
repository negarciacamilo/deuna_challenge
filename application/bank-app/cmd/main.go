package main

import (
	"github.com/negarciacamilo/deuna_challenge/application/bank-app/http"
	"github.com/spf13/viper"
	"log"
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

	if err = http.NewRouter().Run(viper.GetString("BANK_PORT")); err != nil {
		log.Panic(err)
	}
}
