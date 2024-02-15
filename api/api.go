package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"weather-forecaster/forecaster"
)

const (
	addr   = "localhost:3000"
	apiUrl = "/api/{city}"
)

type Api struct {
	mux *http.ServeMux
	fc  forecaster.Forecaster
}

func MakeApi(f forecaster.Forecaster) *Api {
	return &Api{mux: http.DefaultServeMux, fc: f}
}
func (a *Api) cityForecast(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	city := req.PathValue("city")

	all, err := io.ReadAll(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Errorf("cannot read data; %w", err).Error()))
		return
	}

	var p forecaster.Params
	err = json.Unmarshal(all, &p)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(fmt.Errorf("cannot decode request body; %w", err).Error()))
		return
	}

	if !p.Temperature && !p.Precipitation {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("please choose forecast"))
		return
	}
	if p.Days == 0 {
		p.Days = 1
	}

	forecast, err := a.fc.GetCityForecast(city, p)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Errorf("cannot find forecast for given params; %w", err).Error()))
		return
	}

	data, err := json.Marshal(forecast)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Errorf("cannot encode data; %w", err).Error()))
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}

func (a *Api) Serve() error {
	log.Printf("Starting API on address %v/api/", addr)
	a.mux.HandleFunc(apiUrl, a.cityForecast)
	return http.ListenAndServe(addr, a.mux)
}
