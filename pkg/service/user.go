package service

import (
	"ChatServer/pkg/repository"
	"context"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, nickname string) (*int, error) {
	return s.repo.Register(ctx, nickname)
}
