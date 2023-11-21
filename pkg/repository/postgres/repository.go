package postgres

import (
	"ChatServer/pkg/repository/postgres/internal"
	"ChatServer/pkg/types"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Chat interface {
	CreateChat(ctx context.Context, users []int) (*int, error)
	DeleteChat(ctx context.Context, chatId int) error
	GetMissedMessagesFromChats(ctx context.Context, userTo int, chatIds []int, lastSeen time.Time) ([]types.WsMessageOut, error)
	GetAllChats(ctx context.Context, userId int) ([]int, error)
	GetChatMembers(ctx context.Context, chatId int) ([]int, error)
}

type Message interface {
	AddMessage(ctx context.Context, fromId, chatId int, text string) error
}

type User interface {
	Register(ctx context.Context, nickname string) (*int, error)
}

type Info interface {
	GetUsersByPrefix(ctx context.Context, prefix string) (types.UsersList, error)
}

type Device interface {
	RegisterDevice(ctx context.Context, deviceId string, userId int) error
	SaveLastActiveTime(ctx context.Context, userId int, deviceId string, lastTime time.Time) error
	GetLastActiveTime(ctx context.Context, deviceId string, userId int) (time.Time, error)
}

type Repository struct {
	Chat
	Message
	User
	Info
	Device
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Chat:    internal.NewChatPostgres(db),
		User:    internal.NewUserPostgres(db),
		Info:    internal.NewInfoPostgres(db),
		Message: internal.NewMessagePostgres(db),
		Device:  internal.NewDevicePostgres(db),
	}
}
