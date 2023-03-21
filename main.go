package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nemo984/weatherapp/external/airquality"
	weatherAPI "github.com/nemo984/weatherapp/external/weather"

	"github.com/nemo984/weatherapp/middleware"
	"github.com/nemo984/weatherapp/user"
	"github.com/nemo984/weatherapp/weather"
)

func main() {
	r := gin.Default()
	ws := weather.NewService(&airquality.AirQualityAPI{}, &weatherAPI.WeatherAPI{})
	wh := weather.NewHandler(ws)

	us := user.NewService()

	r.Use(gin.Recovery())
	r.Use(middleware.AttachLocation())
	r.GET("/weather", wh.GetWeather)
	r.GET("/air-quality", wh.GetAirQuality)

	usersRoute := r.Group("/users")
	usersRoute.Use(middleware.AttachUser(us))
	{

		remindersRoute := r.Group("/reminders")
		{
			remindersRoute.GET("", wh.GetAirQuality)
			reminderSettingsRoute := remindersRoute.Group("/settings")
			{
				reminderSettingsRoute.GET("", wh.GetAirQuality)
				reminderSettingsRoute.POST("", wh.GetAirQuality)
				reminderSettingsRoute.PUT("/:reminder-id", wh.GetAirQuality)
			}
		}
	}

	r.Run()
}

func initRoutes(r *gin.Engine) {

}
