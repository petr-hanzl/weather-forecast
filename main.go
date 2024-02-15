package main

import (
	"log"
	"net/http"
	"weather-forecaster/api"
	"weather-forecaster/forecaster"
)

func main() {
	f := forecaster.MakeForecaster(http.DefaultClient)
	a := api.MakeApi(f)
	if err := a.Serve(); err != nil {
		log.Fatal(err)
	}
}
