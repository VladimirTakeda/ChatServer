package service

import (
	"ChatServer/pkg/repository/postgres"
	"ChatServer/pkg/types"
	"context"
)

type InfoService struct {
	repo postgres.Info
}

func NewInfoService(repo postgres.Info) *InfoService {
	return &InfoService{repo: repo}
}

func (s *InfoService) GetUsersByPrefix(ctx context.Context, prefix string) (types.UsersList, error) {
	return s.repo.GetUsersByPrefix(ctx, prefix)
}
