package reminder

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReminderType string

const (
	Weather    ReminderType = "WEATHER"
	AirQuality ReminderType = "AIR_QUALITY"
	// RainingForecast ReminderType = "Raining Forecast"
)

func (rt ReminderType) String() string {
	switch rt {
	case Weather:
		return "Weather"
	case AirQuality:
		return "Air Quality"
	default:
		return "Unknown"
	}
}

type ReminderOption string

const (
	Periodic  ReminderOption = "PERIODIC"
	TimeOfDay ReminderOption = "TIME_OF_DAY"
)

func (ro ReminderOption) String() string {
	switch ro {
	case Periodic:
		return "Perodic"
	case TimeOfDay:
		return "Time of the Day"
	default:
		return "Unknown"
	}
}

type Reminder struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserDeviceID     string             `bson:"user_device_id" json:"userDeviceId"`
	Type             ReminderType       `bson:"type" json:"type" binding:"required"`
	Option           ReminderOption     `bson:"option" json:"option" binding:"required"`
	LastRemindedTime time.Time          `bson:"last_reminded_time" json:"lastRemindedTime"`
	RemindAgainOn    time.Time          `bson:"remind_again_on" json:"remindAgainOn"`
	PeriodicDuration Duration           `bson:"periodic_duration" json:"periodicDuration,omitempty"`
	TimeOfDay        TimeHours          `bson:"time_of_day" json:"timeOfDay,omitempty"`
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type TimeHours struct {
	time.Time
}

func (th *TimeHours) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	parsedTime, err := time.ParseInLocation("15:04", s, time.UTC)
	if err != nil {
		return err
	}

	th.Time = parsedTime.In(time.Local)
	return nil
}

func (th TimeHours) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(th.Format("15:04"))), nil
}
