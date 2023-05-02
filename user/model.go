package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DeviceID string             `bson:"device_id"`
	Location Location           `bson:"location"`
	FCMToken string             `bson:"fcm_token" json:"fcmToken"`
}

type Location struct {
	Lat, Lon float64
}

type UpdateFCMTokenRequest struct {
	FCMToken string `json:"fcmToken"`
}
