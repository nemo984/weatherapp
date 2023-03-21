package reminder

type ReminderRepositary struct {
}

// retrieves reminder that are close to to remind time
func (repo *ReminderRepositary) GetRemindersToRemind() []*Reminder {
	return []*Reminder{}
}

func (repo *ReminderRepositary) UpdateReminder(reminder *Reminder) error {
	return nil
}
