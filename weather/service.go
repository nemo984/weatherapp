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
)

type IService interface {
	GetCurrentWeather(ctx context.Context, req WeatherRequest) (*weather.WeatherResponse, error)
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

func (s *Service) GetCurrentWeather(ctx context.Context, req WeatherRequest) (*weather.WeatherResponse, error) {
	fmt.Println(s.cache.Items())
	weatherKey := getCacheKeyWithLocation(req.Lat, req.Lon)
	w, found := s.cache.Get(weatherKey)
	if found {
		return w.(*weather.WeatherResponse), nil
	}

	weather, err := s.weatherAPI.GetWeather(ctx, req.Lat, req.Lon)
	if err != nil {
		return weather, err
	}
	s.cache.Set(weatherKey, weather, 0)
	return weather, nil
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

func getCacheKeyWithLocation(lat, lon float64) string {
	return fmt.Sprintf(weatherCacheKey, geohash.EncodeWithPrecision(lat, lon, 6))
}
