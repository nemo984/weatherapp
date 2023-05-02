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

func NewService(repo ReminderRepositary) *reminderService {
	return &reminderService{reminderRepo: repo}
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
				reminderOption := newReminderOptionStrategy(&reminder)
				if reminderOption.ShouldRemind() {
					fmt.Printf("reminded: %+v\n", reminder)
					rs.remind(ctx, reminder)
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

	noti := &fcm.Notification{}
	noti.Title = fmt.Sprintf("%s %s Reminder", reminder.Option, reminder.Type)
	if reminder.Type == Weather {
		wtr, err := rs.weatherService.GetCurrentWeather(ctx, weather.WeatherRequest{
			Lat: user.Location.Lat,
			Lon: user.Location.Lon,
		})
		if err != nil {
			return err
		}

		mainTemp := wtr.Current.Main
		noti.Body = fmt.Sprintf("Temp: %f (%s), Feels like: %f, High %f, Low %f", mainTemp.Temp, wtr.Current.Weather[0].Description, mainTemp.FeelsLike, mainTemp.TempMax, mainTemp.TempMin)
	} else {
		airQuality, err := rs.weatherService.GetAirQuality(ctx)
		if err != nil {
			return err
		}
		noti.Body = fmt.Sprintf("Current PM 2.5: %f", airQuality.Data.Iaqi.Pm25.V)
	}
	rs.SendPushNotification(context.Background(), "Send test", noti)
	return nil
}

func (rs *reminderService) SendPushNotification(ctx context.Context, deviceToken string, noti *fcm.Notification) error {
	// Create the message to be sent.
	msg := &fcm.Message{
		To:           deviceToken,
		Notification: noti,
	}

	token := "AAAAUscJt1A:APA91bEOdp_yKsWG4xJVmrMHMDZ21nXiGsTW7QqxVIEoO_xf9N309saV1xCI6_2wmNKSZS680cSynEygheyTNEvaJQnU1KyCTeszBf8gEN258gyVksYeSS_AMUWWIu_WRUNfbq_Qwpgo"
	client, err := fcm.NewClient(token)
	if err != nil {
		log.Fatalln(err)
	}
	rs.fcmClient = client

	response, err := rs.fcmClient.SendWithContext(ctx, msg)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", response)
	return nil
}
