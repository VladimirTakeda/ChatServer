package service

import (
	"ChatServer/pkg/repository"
	"context"
)

type ChatService struct {
	repo repository.Chat
}

func NewChatService(repo repository.Chat) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) CreateChat(ctx context.Context, couriers []int) (*int, error) {
	return s.repo.CreateChat(ctx, couriers)
}

func (s *ChatService) DeleteChat(ctx context.Context, chatId int) error {
	return s.repo.DeleteChat(ctx, chatId)
}
