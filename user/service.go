package user

import (
	"context"
)

type Service interface {
	GetOrCreateUser(ctx context.Context, deviceID string) (User, error)
	UpdateUser(ctx context.Context, user *User) error
}

type service struct {
	userRepo UserRepositary
}

func NewService(repo UserRepositary) *service {
	return &service{
		userRepo: repo,
	}
}

func (s *service) GetOrCreateUser(ctx context.Context, deviceID string) (User, error) {
	return s.userRepo.UpsertUser(ctx, &User{
		DeviceID: deviceID,
	})
}

func (s *service) UpdateUser(ctx context.Context, user *User) error {
	_, err := s.userRepo.UpsertUser(ctx, user)
	return err
}
