package main

import (
	"context"
	"log"
	"time"

	"github.com/appleboy/go-fcm"
	"github.com/gin-gonic/gin"
	"github.com/nemo984/weatherapp/external/airquality"
	weatherAPI "github.com/nemo984/weatherapp/external/weather"
	"github.com/nemo984/weatherapp/reminder"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nemo984/weatherapp/middleware"
	"github.com/nemo984/weatherapp/user"
	"github.com/nemo984/weatherapp/weather"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	token := "AAAAUscJt1A:APA91bEOdp_yKsWG4xJVmrMHMDZ21nXiGsTW7QqxVIEoO_xf9N309saV1xCI6_2wmNKSZS680cSynEygheyTNEvaJQnU1KyCTeszBf8gEN258gyVksYeSS_AMUWWIu_WRUNfbq_Qwpgo"
	fClient, err := fcm.NewClient(token)
	if err != nil {
		log.Fatalln(err)
	}

	userRepo := user.NewRepositary(client.Database("weatherapp"))
	us := user.NewService(*userRepo)
	ws := weather.NewService(&airquality.AirQualityAPI{}, &weatherAPI.WeatherAPI{})

	reminderRepo := reminder.NewRepositary(client.Database("weatherapp"))
	reminderService := reminder.NewService(*reminderRepo, us, ws, fClient)
	reminderHandler := reminder.Handler{ReminderService: reminderService}

	r := gin.Default()
	wh := weather.NewHandler(ws)

	userHandler := user.NewHandler(us)

	r.Use(gin.Recovery())
	r.Use(middleware.AttachLocation())
	r.GET("/weather", wh.GetWeather)
	r.GET("/air-quality", wh.GetAirQuality)

	usersRoute := r.Group("/users")
	usersRoute.PUT("/fcm-token", userHandler.UpdateFCMToken)
	usersRoute.Use(middleware.AttachUser(us))
	{

		remindersRoute := usersRoute.Group("/reminders")
		{
			remindersRoute.GET("", reminderHandler.GetReminders)
			remindersRoute.POST("", reminderHandler.UpsertReminder)
			remindersRoute.PUT("", reminderHandler.UpsertReminder)
		}
	}

	go reminderService.StartReminderJob(context.Background())
	r.Run()
}

func initRoutes(r *gin.Engine) {

}
