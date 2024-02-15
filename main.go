package main

import (
	"net/http"
	"weather-forecaster/forecaster"
)

func main() {
	f := forecaster.MakeForecaster(http.DefaultClient)
	f.GetCityForecast("bRnO", 2)
}
