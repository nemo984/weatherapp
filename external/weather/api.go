package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	appID            = "b095e2533215e8ddc7cbf38e41b416d1"
	weatherEndpoint  = "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric&lang=TH"
	forecastEndpoint = "https://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&appid=%s&units=metric&lang=TH"
)

type IWeatherAPI interface {
	GetWeather(ctx context.Context, lat, lon float64) (*WeatherResponse, error)
	GetForecast(ctx context.Context, lat, lon float64) (*ForecastResponse, error)
}

type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type WeatherAPI struct {
	client http.Client
}

func (w *WeatherAPI) GetWeather(ctx context.Context, lat, lon float64) (*WeatherResponse, error) {
	res := &WeatherResponse{}
	url := fmt.Sprintf(weatherEndpoint, lat, lon, appID)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	r, err := w.client.Do(req)
	if err != nil {
		return res, err
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(res); err != nil {
		return res, err
	}
	return res, nil
}

type ForecastResponse struct {
	Cod     string     `json:"cod"`
	Message int        `json:"message"`
	Cnt     int        `json:"cnt"`
	List    []Forecast `json:"list"`
	City    struct {
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

type Forecast struct {
	Dt   int `json:"dt"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
		Humidity  int     `json:"humidity"`
		TempKf    float64 `json:"temp_kf"`
	} `json:"main"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Visibility int     `json:"visibility"`
	Pop        float64 `json:"pop"`
	Sys        struct {
		Pod string `json:"pod"`
	} `json:"sys"`
	DtTxt string `json:"dt_txt"`
}

func (w *WeatherAPI) GetForecast(ctx context.Context, lat, lon float64) (*ForecastResponse, error) {
	res := &ForecastResponse{}
	url := fmt.Sprintf(forecastEndpoint, lat, lon, appID)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	r, err := w.client.Do(req)
	if err != nil {
		return res, err
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(res); err != nil {
		return res, err
	}
	return res, nil
}
