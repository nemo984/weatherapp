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
	GetUserReminders(ctx context.Context, deviceID string) ([]Reminder, error)
	StartReminderJob(ctx context.Context) error
	UpsertReminder(ctx context.Context, reminder Reminder) (Reminder, error)
}

type reminderService struct {
	fcmClient      *fcm.Client
	weatherService weather.IService
	userService    user.Service
	reminderRepo   ReminderRepositary
}

func NewService(repo ReminderRepositary, userService user.Service, weatherService weather.IService, fcmClient *fcm.Client) *reminderService {
	return &reminderService{reminderRepo: repo, userService: userService, weatherService: weatherService, fcmClient: fcmClient}
}

func (rs *reminderService) GetUserReminders(ctx context.Context, deviceID string) ([]Reminder, error) {
	return rs.reminderRepo.ListUserReminders(ctx, deviceID)
}

func (rs *reminderService) UpsertReminder(ctx context.Context, reminder Reminder) (Reminder, error) {
	reminderStrategy := newReminderOptionStrategy(&reminder)
	reminderStrategy.CalculateRemindAgainOn()
	return rs.reminderRepo.Upsert(ctx, reminder)
}

func (rs *reminderService) StartReminderJob(ctx context.Context) error {
	ticker := time.NewTicker(reminderJobInterval)
	for {
		select {
		case <-ticker.C:
			log.Println("ticked, reminder job loop")
			reminders, err := rs.reminderRepo.GetRemindersToRemind(ctx)
			if err != nil {
				log.Printf("error reminder job, get reminders: %w\n", err)
			}
			log.Printf("got %v reminders to remind\n", len(reminders))

			for _, reminder := range reminders {
				reminderOption := newReminderOptionStrategy(reminder)
				if reminderOption.ShouldRemind() {
					fmt.Printf("reminded: %+v\n", reminder)
					rs.remind(ctx, *reminder)
					reminderOption.CalculateRemindAgainOn()
				}
			}

			err = rs.reminderRepo.UpsertMany(ctx, reminders)
			if err != nil {
				log.Printf("error reminder job, upsert many: %w\n", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (rs *reminderService) remind(ctx context.Context, reminder Reminder) error {
	user, err := rs.userService.GetOrCreateUser(ctx, reminder.UserDeviceID)
	if err != nil {
		return err
	}

	lat, lon := 13.7563, 100.5018
	if user.Location.Lon != 0 {
		lat = user.Location.Lon
	}
	if user.Location.Lat != 0 {
		lon = user.Location.Lat
	}

	noti := &fcm.Notification{}
	noti.Title = fmt.Sprintf("%s %s Reminder", reminder.Option, reminder.Type)
	if reminder.Type == Weather {
		wtr, err := rs.weatherService.GetCurrentWeather(ctx, weather.WeatherRequest{
			Lat: lat,
			Lon: lon,
		})
		if err != nil {
			return err
		}

		mainTemp := wtr.Current.Main
		noti.Body = fmt.Sprintf("Temp: %.2f (%s), Feels like: %.2f, High %.2f, Low %.2f", mainTemp.Temp, wtr.Current.Weather[0].Description, mainTemp.FeelsLike, mainTemp.TempMax, mainTemp.TempMin)
	} else {
		airQuality, err := rs.weatherService.GetAirQuality(ctx)
		if err != nil {
			return err
		}
		noti.Body = fmt.Sprintf("Current PM 2.5: %.2f", airQuality.Data.Iaqi.Pm25.V)
	}
	rs.SendPushNotification(context.Background(), user.FCMToken, noti)
	return nil
}

func (rs *reminderService) SendPushNotification(ctx context.Context, deviceToken string, noti *fcm.Notification) error {
	// Create the message to be sent.
	msg := &fcm.Message{
		To:           deviceToken,
		Notification: noti,
	}

	response, err := rs.fcmClient.SendWithContext(ctx, msg)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", response)
	return nil
}
