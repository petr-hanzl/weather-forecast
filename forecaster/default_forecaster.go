package forecaster

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	tempParam          = "temperature_180m"
	precipitationParam = "Precipitation"
)

func MakeForecaster(client *http.Client) Forecaster {
	return &defaultForecaster{client: client}
}

type defaultForecaster struct {
	client *http.Client
}

type geoResponse struct {
	Results []City
}

func (f *defaultForecaster) getCity(cityStr string) (City, error) {
	endpoint := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json", url.QueryEscape(cityStr))
	resp, err := f.client.Get(endpoint)
	if err != nil {
		return City{}, fmt.Errorf("cannot send request to geo location api; %w", err)
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return City{}, fmt.Errorf("cannot read information about city; %w", err)
	}

	var geoResp geoResponse
	err = json.Unmarshal(all, &geoResp)
	if err != nil {
		return City{}, fmt.Errorf("cannot decode information about city; %w", err)
	}

	return geoResp.Results[0], nil
}

type forecastResponse struct {
	Tmp temp `json:"hourly"`
}

type temp struct {
	Time        []string
	Temperature []float64 `json:"temperature_180m"`
}

func (f *defaultForecaster) GetCityForecast(cityStr string, params Params) (City, error) {
	city, err := f.getCity(cityStr)
	if err != nil {
		return city, err
	}
	var apiParams string
	switch {
	case params.Temperature && params.Precipitation:
		apiParams = tempParam + "," + precipitationParam
	case params.Temperature:
		apiParams = tempParam
	default:
		apiParams = precipitationParam
	}
	endpoint := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.6f&longitude=%.6f&hourly=%v&forecast_days=%v", city.Longitude, city.Latitude, apiParams, params.Days)
	resp, err := f.client.Get(endpoint)
	if err != nil {
		return City{}, fmt.Errorf("cannot send request to wather api; %w", err)
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return City{}, fmt.Errorf("cannot read weather response; %w", err)
	}

	var fc forecastResponse
	err = json.Unmarshal(all, &fc)
	if err != nil {
		return city, fmt.Errorf("cannot unmarshal Forecasts; %w", err)
	}

	// fill data into city
	for i, _ := range fc.Tmp.Time {
		city.Forecasts = append(city.Forecasts, Forecast{
			Time: fc.Tmp.Time[i],
			Temp: fc.Tmp.Temperature[i],
		})
	}

	return city, nil
}
