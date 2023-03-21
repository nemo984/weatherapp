package user

import (
	"context"
)

type Service interface {
	GetOrCreateUser(ctx context.Context, deviceID string) *User
}

type service struct {
	userRepo UserRepositary
}

func NewService() *service {
	return &service{}
}

func (s *service) GetOrCreateUser(ctx context.Context, deviceID string) *User {
	return &User{
		DeviceID: deviceID,
	}
}
