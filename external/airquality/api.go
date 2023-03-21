package airquality

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type IAirQualityAPI interface {
	GetGeoLocalizedFeed(ctx context.Context) (*AirQualityResponse, error)
}

type AirQualityResponse struct {
	Status string `json:"status"`
	Data   struct {
		Aqi          int `json:"aqi"`
		Idx          int `json:"idx"`
		Attributions []struct {
			URL  string `json:"url"`
			Name string `json:"name"`
			Logo string `json:"logo,omitempty"`
		} `json:"attributions"`
		City struct {
			Geo      []float64 `json:"geo"`
			Name     string    `json:"name"`
			URL      string    `json:"url"`
			Location string    `json:"location"`
		} `json:"city"`
		Dominentpol string `json:"dominentpol"`
		Iaqi        struct {
			H struct {
				V int `json:"v"`
			} `json:"h"`
			No2 struct {
				V float64 `json:"v"`
			} `json:"no2"`
			O3 struct {
				V float64 `json:"v"`
			} `json:"o3"`
			P struct {
				V float64 `json:"v"`
			} `json:"p"`
			Pm10 struct {
				V int `json:"v"`
			} `json:"pm10"`
			Pm25 struct {
				V int `json:"v"`
			} `json:"pm25"`
			T struct {
				V float64 `json:"v"`
			} `json:"t"`
			W struct {
				V float64 `json:"v"`
			} `json:"w"`
			Wg struct {
				V int `json:"v"`
			} `json:"wg"`
		} `json:"iaqi"`
		Time struct {
			S   string    `json:"s"`
			Tz  string    `json:"tz"`
			V   int       `json:"v"`
			Iso time.Time `json:"iso"`
		} `json:"time"`
		Forecast struct {
			Daily struct {
				O3 []struct {
					Avg int    `json:"avg"`
					Day string `json:"day"`
					Max int    `json:"max"`
					Min int    `json:"min"`
				} `json:"o3"`
				Pm10 []struct {
					Avg int    `json:"avg"`
					Day string `json:"day"`
					Max int    `json:"max"`
					Min int    `json:"min"`
				} `json:"pm10"`
				Pm25 []struct {
					Avg int    `json:"avg"`
					Day string `json:"day"`
					Max int    `json:"max"`
					Min int    `json:"min"`
				} `json:"pm25"`
			} `json:"daily"`
		} `json:"forecast"`
		Debug struct {
			Sync time.Time `json:"sync"`
		} `json:"debug"`
	} `json:"data"`
}

type AirQualityAPI struct {
	client http.Client
}

func (aq *AirQualityAPI) GetGeoLocalizedFeed(ctx context.Context) (*AirQualityResponse, error) {
	res := &AirQualityResponse{}
	url := "https://api.waqi.info/feed/here/?token=06ea7e42dd59339aa1a59396a8021d29f5a3f212"
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	r, err := aq.client.Do(req)
	if err != nil {
		return res, err
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(res); err != nil {
		return res, err
	}
	return res, nil
}
