package user

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		userService: s,
	}
}

func (h Handler) UpdateFCMToken(c *gin.Context) {
	req := UpdateFCMTokenRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	latFloat, _ := strconv.ParseFloat(strings.TrimSpace(c.GetHeader("X-Lat")), 64)
	lonFloat, _ := strconv.ParseFloat(strings.TrimSpace(c.GetHeader("X-Lon")), 64)
	user := &User{
		DeviceID: c.GetHeader("X-Device-ID"),
		Location: Location{
			Lat: latFloat,
			Lon: lonFloat,
		},
		FCMToken: req.FCMToken,
	}

	err := h.userService.UpdateUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
