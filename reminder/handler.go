package reminder

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nemo984/weatherapp/middleware"
)

type Handler struct {
	ReminderService ReminderService
}

func (h *Handler) GetReminders(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	reminders, err := h.ReminderService.GetUserReminders(c, user.DeviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *Handler) UpsertReminder(c *gin.Context) {
	user := middleware.GetUserFromContext(c)
	reminder := Reminder{}
	if err := c.BindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	reminder.UserDeviceID = user.DeviceID

	reminder, err := h.ReminderService.UpsertReminder(c, reminder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, reminder)
}
