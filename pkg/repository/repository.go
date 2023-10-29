package repository

import (
	"ChatServer/pkg/types"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Chat interface {
	CreateChat(ctx context.Context, users []int) (*int, error)
	DeleteChat(ctx context.Context, chatId int) error
}

type Message interface {
	Create(ctx context.Context, fromId, chatId int) error
}

type User interface {
	Register(ctx context.Context, nickname string) (*int, error)
}

type Info interface {
	GetUsersByPrefix(ctx context.Context, prefix string) (types.UsersList, error)
}

type Repository struct {
	Chat
	Message
	User
	Info
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Chat: NewChatPostgres(db),
		User: NewUserPostgres(db),
		Info: NewInfoPostgres(db),
	}
}
