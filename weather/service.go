package weather

import (
	"context"
	"fmt"
	"time"

	"github.com/mmcloughlin/geohash"
	"github.com/nemo984/weatherapp/external/airquality"
	"github.com/nemo984/weatherapp/external/weather"
	"github.com/patrickmn/go-cache"
)

const (
	weatherCacheKey    = "weatherKey-%s"
	airQualityCacheKey = "airQualityKey"
	geohashPrecision   = 6
)

type IService interface {
	GetCurrentWeather(ctx context.Context, req WeatherRequest) (*WeatherResponse, error)
	GetAirQuality(ctx context.Context) (*airquality.AirQualityResponse, error)
}

type Service struct {
	weatherAPI    weather.IWeatherAPI
	airQualityAPI airquality.IAirQualityAPI

	cache *cache.Cache
}

var _ IService = (*Service)(nil)

func NewService(airQualityAPI airquality.IAirQualityAPI, weatherAPI weather.IWeatherAPI) *Service {
	return &Service{
		airQualityAPI: airQualityAPI,
		weatherAPI:    weatherAPI,
		cache:         cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (s *Service) GetCurrentWeather(ctx context.Context, req WeatherRequest) (*WeatherResponse, error) {
	weatherKey := getCacheKeyWithLocation(weatherCacheKey, req.Lat, req.Lon)
	w, found := s.cache.Get(weatherKey)
	if found {
		return w.(*WeatherResponse), nil
	}

	we, err := s.weatherAPI.GetWeather(ctx, req.Lat, req.Lon)
	if err != nil {
		return &WeatherResponse{}, err
	}

	res := &WeatherResponse{
		Current: *we,
	}
	forecast, err := s.weatherAPI.GetForecast(ctx, req.Lat, req.Lon)
	if err != nil {
		return res, err
	}

	res.City = forecast.City
	dayGroupMap := make(map[string][]ForecastHour)
	for _, forecast := range forecast.List {
		t := time.Unix(int64(forecast.Dt), 0)
		weekDay := t.Weekday().String()
		hour := fmt.Sprintf("%.2f", float64(t.Hour()))
		dayGroupMap[weekDay] = append(dayGroupMap[weekDay], ForecastHour{
			Forecast: forecast,
			Hour:     hour,
		})
	}
	for day, forecasts := range dayGroupMap {
		f := ForecastDay{
			Day: day,
		}
		for _, forecast := range forecasts {
			f.Hours = append(f.Hours, forecast)
		}
		res.Forecasts = append(res.Forecasts, f)
	}

	s.cache.Set(weatherKey, res, 0)
	return res, nil
}

func (s *Service) GetAirQuality(ctx context.Context) (*airquality.AirQualityResponse, error) {
	aq, found := s.cache.Get(airQualityCacheKey)
	if found {
		return aq.(*airquality.AirQualityResponse), nil
	}

	airQuality, err := s.airQualityAPI.GetGeoLocalizedFeed(ctx)
	if err != nil {
		return &airquality.AirQualityResponse{}, err
	}
	s.cache.Set(airQualityCacheKey, airQuality, 0)
	return airQuality, nil
}

func getCacheKeyWithLocation(key string, lat, lon float64) string {
	return fmt.Sprintf(key, geohash.EncodeWithPrecision(lat, lon, geohashPrecision))
}
