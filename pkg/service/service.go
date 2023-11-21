package service

import (
	"ChatServer/pkg/repository/postgres"
	"ChatServer/pkg/repository/redis"
	"ChatServer/pkg/types"
	"context"
	"time"
)

type Chat interface {
	GetMissedMessages(ctx context.Context, userId int, deviceId string) ([]types.WsMessageOut, error)
	CreateChat(ctx context.Context, users []int) (*int, error)
	DeleteChat(ctx context.Context, chatId int) error
	GetChatMembers(ctx context.Context, chatId int) ([]int, error)
}

type Device interface {
	RegisterDevice(ctx context.Context, deviceId string, userId int) error
	SaveLastActiveTime(ctx context.Context, userId int, deviceId string, lastTime time.Time) error
	GetLastActiveTime(ctx context.Context, deviceId string, userId int) (time.Time, error)
}

type User interface {
	Register(ctx context.Context, nickname string) (*int, error)
}

type Info interface {
	GetUsersByPrefix(ctx context.Context, prefix string) (types.UsersList, error)
}

type Message interface {
	AddMessage(ctx context.Context, fromId, chatId int, text string) error
}

type Service struct {
	Chat
	User
	Info
	Message
	Device
}

func NewService(repos *postgres.Repository, cache *redis.Cache) *Service {
	return &Service{
		Chat:    NewChatService(repos.Chat, repos.Device, cache.Chat),
		User:    NewUserService(repos.User),
		Info:    NewInfoService(repos.Info),
		Message: NewMessageService(repos.Message),
		Device:  NewDeviceService(repos.Device),
	}
}
