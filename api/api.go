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

type errorResponse struct {
	Code    int
	Message string
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
		e := errorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("cannot read data; %v", err),
		}
		marshal, err := json.Marshal(e)
		if err != nil {
			rw.Write([]byte("error"))
			return
		}

		rw.Write(marshal)
		return
	}

	var p forecaster.Params
	err = json.Unmarshal(all, &p)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		e := errorResponse{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("cannot decode request body; %v", err),
		}
		marshal, err := json.Marshal(e)
		if err != nil {
			rw.Write([]byte("error"))
			return
		}

		rw.Write(marshal)
		return
	}

	if !p.Temperature && !p.Precipitation {
		rw.WriteHeader(http.StatusBadRequest)
		e := errorResponse{
			Code:    http.StatusBadRequest,
			Message: "please choose forecast",
		}
		marshal, err := json.Marshal(e)
		if err != nil {
			rw.Write([]byte("error"))
			return
		}

		rw.Write(marshal)
		return
	}
	if p.Days == 0 {
		p.Days = 1
	}

	forecast, err := a.fc.GetCityForecast(city, p)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		e := errorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("cannot find forecast for given params; %v", err),
		}
		marshal, err := json.Marshal(e)
		if err != nil {
			rw.Write([]byte("error"))
			return
		}

		rw.Write(marshal)
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
