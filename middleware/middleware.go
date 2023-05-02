package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nemo984/weatherapp/user"
)

const (
	ctxUserKey    = "userKey"
	ctxUserLatKey = "userLatKey"
	ctxUserLonKey = "userLonKey"
)

func AttachUser(userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.GetHeader("X-Device-ID")
		if deviceID == "" {
			c.JSON(http.StatusBadRequest, "missing X-Device-ID header")
			return
		}
		user, err := userService.GetOrCreateUser(c, deviceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Set(ctxUserKey, user)
	}
}

func AttachLocation() gin.HandlerFunc {
	return func(c *gin.Context) {
		lat := c.Query("lat")
		lon := c.Query("lon")
		if lat == "" || lon == "" {
			c.Next()
			return
		}
		latFloat, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		lonFloat, err := strconv.ParseFloat(lon, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.Set(ctxUserLatKey, latFloat)
		c.Set(ctxUserLonKey, lonFloat)
	}
}

func GetUserFromContext(c *gin.Context) user.User {
	value, _ := c.Get(ctxUserKey)
	return value.(user.User)
}

func GetLocationFromContext(c *gin.Context) (lat float64, lon float64) {
	return c.GetFloat64(ctxUserLatKey), c.GetFloat64(ctxUserLonKey)
}
