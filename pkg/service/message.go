package service

import (
	"ChatServer/pkg/repository/postgres"
	"context"
)

type MessageService struct {
	repo postgres.Message
}

func NewMessageService(repo postgres.Message) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) AddMessage(ctx context.Context, fromId, chatId int, text string, attachments []string) error {
	return s.repo.AddMessage(ctx, fromId, chatId, text, attachments)
}
