package weather

import "github.com/nemo984/weatherapp/external/weather"

type WeatherResponse struct {
	Current   weather.WeatherResponse `json:"current"`
	Forecasts []ForecastDay           `json:"forecasts"`
	City      struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}

type ForecastDay struct {
	Day   string         `json:"day"` // Monday, Tuesday, ...
	Hours []ForecastHour `json:"hours"`
}

type ForecastHour struct {
	Forecast weather.Forecast `json:"forecast"`
	Hour     string           `json:"hour"`
}

type WeatherRequest struct {
	Lat, Lon float64
}

type AirQuality struct{}
