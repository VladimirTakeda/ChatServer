package service

import (
	"ChatServer/pkg/repository"
)

type Service struct {
	repository.Chat
	repository.User
	repository.Info
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Chat: NewChatService(repos.Chat),
		User: NewUserService(repos.User),
		Info: NewInfoService(repos.Info),
	}
}
