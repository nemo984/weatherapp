package reminder

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReminderType string

const (
	Weather    ReminderType = "WEATHER"
	AirQuality ReminderType = "AIR_QUALITY"
	// RainingForecast ReminderType = "Raining Forecast"
)

type ReminderOption int

const (
	Periodic ReminderOption = iota
	TimeOfDay
)

type Reminder struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	UserDeviceID     string             `bson:"user_device_id"`
	Type             ReminderType       `bson:"type"`
	Option           ReminderOption     `bson:"option"`
	LastRemindedTime time.Time          `bson:"last_reminded_time"`
	RemindAgainOn    time.Time          `bson:"remind_again_on"`
	PeriodicDuration time.Duration      `bson:"periodic_duration"`
	TimeOfDay        time.Time          `bson:"time_of_day"`
}
