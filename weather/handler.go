package weather

import (
	"fmt"
	"net/http"

	"github.com/nemo984/weatherapp/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service IService
}

func NewHandler(service IService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetWeather(c *gin.Context) {
	lat, lon := middleware.GetLocationFromContext(c)
	weather, _ := h.service.GetCurrentWeather(c, WeatherRequest{
		Lat: lat,
		Lon: lon,
	})
	c.JSON(http.StatusOK, weather)
}

func (h *Handler) GetAirQuality(c *gin.Context) {
	aq, _ := h.service.GetAirQuality(c)
	fmt.Printf("%+v", aq)
}
