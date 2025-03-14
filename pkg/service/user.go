package service

import (
	"ChatServer/pkg/repository/postgres"
	"context"
)

type UserService struct {
	repo postgres.User
}

func NewUserService(repo postgres.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, nickname string) (*int, error) {
	userID, err := s.repo.Register(ctx, nickname)
	if err != nil {
		return nil, err
	}
	return userID, nil
}
