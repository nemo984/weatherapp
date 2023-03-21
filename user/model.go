package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	DeviceID        string
	Location        Location
	ReminderEnabled bool
}

type Location struct {
	Lat, Lon float64
}

// Collections
//
