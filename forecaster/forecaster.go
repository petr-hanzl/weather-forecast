package forecaster

type Forecaster interface {
	GetCityForecast(city string, params Params) (City, error)
}

type Forecast struct {
	Time          string
	Temp          float64
	Precipitation float64
}
type City struct {
	Name      string
	Latitude  float64
	Longitude float64
	Timezone  string
	Country   string
	// maximum, 24 hours for hourly results
	Forecasts []Forecast
}

type Params struct {
	Days          int
	Temperature   bool
	Precipitation bool
}
