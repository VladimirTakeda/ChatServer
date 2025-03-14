package service

import (
	minio2 "ChatServer/pkg/repository/minio"
	"ChatServer/pkg/repository/postgres"
	"ChatServer/pkg/repository/redis"
	"ChatServer/pkg/types"
	"context"
	"github.com/minio/minio-go/v7"
	redis2 "github.com/redis/go-redis/v9"
	"io"
	"time"
)

type Chat interface {
	GetMissedMessages(ctx context.Context, userId int, deviceId string) ([]types.WsMessageWithTime, error)
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
	AddMessage(ctx context.Context, message types.WsMessageWithTime, userIDs []int) error
	CreateSubscription(ctx context.Context, topicName string) (*redis2.PubSub, error)
}

type FileManager interface {
	SaveFile(ctx context.Context, fromId, chatId int, file io.Reader, size int64, fileName string) (string, error)
	LoadFile(ctx context.Context, objName string) (*minio.Object, error)
}

type Service struct {
	Chat
	User
	Info
	Message
	Device
	FileManager
}

func NewService(repos *postgres.Repository, cache *redis.Cache, pubsub *redis.PubSub, s3Storage *minio2.S3Storage) *Service {
	return &Service{
		Chat:        NewChatService(repos.Chat, repos.Device, cache.Chat),
		User:        NewUserService(repos.User),
		Info:        NewInfoService(repos.Info),
		Message:     NewMessageService(repos.Message, pubsub.MessageBus),
		Device:      NewDeviceService(repos.Device),
		FileManager: NewFileManagerService(s3Storage),
	}
}

func (s *Service) RegisterUser(ctx context.Context, nickname, deviseHash string) (*int, string, error) {
	userID, err := s.User.Register(ctx, nickname)
	if err != nil {
		return nil, "", err
	}
	err = s.Device.RegisterDevice(ctx, deviseHash, *userID)
	if err != nil {
		return nil, "", err
	}
	return userID, deviseHash, nil
}
