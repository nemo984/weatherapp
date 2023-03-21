package reminder

import "time"

type ReminderOptionStrategy interface {
	ShouldRemind() bool
	Reminded() // set next time to remind on
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
	return time.Since(pr.LastRemindedTime) >= pr.PeriodicDuration
}

func (pr *PeriodicReminder) Reminded() {
	pr.LastRemindedTime = time.Now()
	pr.RemindAgainOn = pr.LastRemindedTime.Add(pr.PeriodicDuration)
}

type TimeOfDayReminder struct {
	*Reminder
}

func (tor *TimeOfDayReminder) ShouldRemind() bool {
	return tor.LastRemindedTime.Before(tor.TimeOfDay) && tor.TimeOfDay.Truncate(time.Minute).Equal(time.Now().Truncate(time.Minute))
}

func (tor *TimeOfDayReminder) Reminded() {
	tor.LastRemindedTime = time.Now()
	tor.RemindAgainOn = tor.TimeOfDay.Add(24 * time.Hour)
	tor.TimeOfDay = tor.RemindAgainOn
}
