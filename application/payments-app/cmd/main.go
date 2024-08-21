package main

import (
	"github.com/negarciacamilo/deuna_challenge/application/environment"
	"github.com/negarciacamilo/deuna_challenge/application/payments-app/http"
	"github.com/spf13/viper"
)

func main() {
	if !environment.IsDockerEnv() {
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
