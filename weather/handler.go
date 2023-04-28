package weather

import (
	"log"
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
	log.Println(lat, lon)
	weather, err := h.service.GetCurrentWeather(c, WeatherRequest{
		Lat: lat,
		Lon: lon,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, weather)
}

func (h *Handler) GetAirQuality(c *gin.Context) {
	aq, err := h.service.GetAirQuality(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, aq)
}
