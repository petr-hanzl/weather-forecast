package forecaster

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Forecaster struct {
	client *http.Client
}

type City struct {
	Name      string
	Latitude  float64
	Longitude float64
	Timezone  string
	Country   string
}

type GeoResponse struct {
	Results []City
}

func (f *Forecaster) GetCoordinates(cityStr string) (City, error) {
	endpoint := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json", url.QueryEscape(cityStr))
	resp, err := f.client.Get(endpoint)
	if err != nil {
		return City{}, fmt.Errorf("cannot make request; %v", err)
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return City{}, fmt.Errorf("cannot read information about city; %v", err)
	}

	var geoResp GeoResponse
	err = json.Unmarshal(all, &geoResp)
	if err != nil {
		return City{}, fmt.Errorf("cannot decode information about city; %v", err)
	}

	return geoResp.Results[0], nil
}

func (f *Forecaster) GetForecast(lat, long float32) {
	endpoint := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.6f&longitude=%.6f&hourly=temperature_2m", 49.2, 16.8)
	resp, err := f.client.Get(endpoint)
	if err != nil {
		log.Fatalf("error making request to Weather API: %w", err)
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Println(string(all))
}
