package reminder

import (
	"time"
)

type ReminderOptionStrategy interface {
	ShouldRemind() bool
	CalculateRemindAgainOn() // set next time to remind on
}

func newReminderOptionStrategy(r *Reminder) ReminderOptionStrategy {
	switch r.Option {
	case Periodic:
		return &PeriodicReminder{r}
	case TimeOfDay:
		return &TimeOfDayReminder{r}
	default:
		return nil
	}
}

type PeriodicReminder struct {
	*Reminder
}

func (pr *PeriodicReminder) ShouldRemind() bool {
	return time.Since(pr.LastRemindedTime) >= pr.PeriodicDuration.Duration
}

func (pr *PeriodicReminder) CalculateRemindAgainOn() {
	pr.LastRemindedTime = time.Now()
	pr.RemindAgainOn = pr.LastRemindedTime.Add(pr.PeriodicDuration.Duration)
}

type TimeOfDayReminder struct {
	*Reminder
}

func (tor *TimeOfDayReminder) ShouldRemind() bool {
	return time.Since(tor.LastRemindedTime) >= 5*time.Second
}

func (tor *TimeOfDayReminder) CalculateRemindAgainOn() {
	tor.LastRemindedTime = time.Now()
	tor.RemindAgainOn = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), tor.TimeOfDay.Hour(), tor.TimeOfDay.Minute(), 0, 0, time.Now().Location())
	if tor.RemindAgainOn.Before(tor.LastRemindedTime) {
		tor.RemindAgainOn.AddDate(0, 0, 1)
	}
}
