package user

type UserRepositary struct{}

func (userRepo *UserRepositary) GetReminderHistories() {}

func (userRepo *UserRepositary) GetRemindersSettings() {

}

func (userRepo *UserRepositary) GetUsersWithReminderEnabled() []User {
	return []User{}
}
