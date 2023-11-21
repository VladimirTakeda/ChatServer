package redis

import (
	"ChatServer/pkg/repository/redis/internal"
	"context"
	"github.com/redis/go-redis/v9"
)

type Chat interface {
	GetChatMembers(ctx context.Context, chatId int) ([]int, error)
	SetChatMembers(ctx context.Context, chatId int, members []int) error
}

type Cache struct {
	Chat
}

func NewCache(db *redis.Client) *Cache {
	return &Cache{
		Chat: internal.NewChatCache(db),
	}
}
