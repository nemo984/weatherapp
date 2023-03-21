package reminder

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appleboy/go-fcm"
	"github.com/nemo984/weatherapp/user"
	"github.com/nemo984/weatherapp/weather"
)

const (
	reminderJobInterval = 5 * time.Second
)

type ReminderService interface {
	StartReminderJob(context context.Context) error
	UpsertReminder(reminder Reminder) *Reminder
	ListReminders(reminder Reminder) []Reminder
}

type reminderService struct {
	fcmClient      *fcm.Client
	weatherService weather.IService
	userService    user.Service
	reminderRepo   ReminderRepositary
}

func NewService() *reminderService {
	return &reminderService{}
}

func (rs *reminderService) StartReminderJob(ctx context.Context) error {
	ticker := time.NewTicker(reminderJobInterval)
	for {
		select {
		case <-ticker.C:
			fmt.Println("ticked")
			reminders := rs.reminderRepo.GetRemindersToRemind()
			for _, reminder := range reminders {
				reminderOption := newReminderOptionStrategy(reminder)
				if reminderOption.ShouldRemind() {
					go func(reminder *Reminder) {
						fmt.Printf("reminded: %+v\n", reminder)
						rs.remind(ctx, reminder)
						reminderOption.Reminded()
						rs.reminderRepo.UpdateReminder(reminder)
					}(reminder)
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (rs *reminderService) remind(ctx context.Context, reminder *Reminder) error {
	user := rs.userService.GetOrCreateUser(ctx, reminder.UserDeviceID)
	if !user.ReminderEnabled {
		return nil
	}
	if reminder.Type == Weather {
		wtr, err := rs.weatherService.GetCurrentWeather(ctx, weather.WeatherRequest{
			Lat: user.Location.Lat,
			Lon: user.Location.Lon,
		})
		if err != nil {
			return err
		}
		fmt.Println(wtr)
	} else {
		airQuality, err := rs.weatherService.GetAirQuality(ctx)
		if err != nil {
			return err
		}
		fmt.Println(airQuality)
	}
	return nil
}

func (rs *reminderService) SendPushNotification(ctx context.Context, deviceToken string) error {
	// Create the message to be sent.
	msg := &fcm.Message{
		To: deviceToken,
		Data: map[string]interface{}{
			"foo": "bar",
		},
		Notification: &fcm.Notification{
			Title: "title",
			Body:  "body",
		},
	}

	// Create a FCM client to send the message.
	token := "AAAAUscJt1A:APA91bEOdp_yKsWG4xJVmrMHMDZ21nXiGsTW7QqxVIEoO_xf9N309saV1xCI6_2wmNKSZS680cSynEygheyTNEvaJQnU1KyCTeszBf8gEN258gyVksYeSS_AMUWWIu_WRUNfbq_Qwpgo"
	client, err := fcm.NewClient(token)
	if err != nil {
		log.Fatalln(err)
	}
	rs.fcmClient = client

	// Send the message and receive the response without retries.
	response, err := rs.fcmClient.SendWithContext(ctx, msg)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", response)
	return nil
}