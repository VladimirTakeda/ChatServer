package service

import (
	"ChatServer/pkg/repository"
	"ChatServer/pkg/types"
	"context"
)

type InfoService struct {
	repo repository.Info
}

func NewInfoService(repo repository.Info) *InfoService {
	return &InfoService{repo: repo}
}

func (s *InfoService) GetUsersByPrefix(ctx context.Context, prefix string) (types.UsersList, error) {
	return s.repo.GetUsersByPrefix(ctx, prefix)
}
