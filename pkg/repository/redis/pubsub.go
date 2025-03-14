package redis

import (
	"ChatServer/pkg/repository/redis/internal"
	"context"
	"github.com/redis/go-redis/v9"
)

// TODO: can be reused for any other Bus RabbitMQ, Cassandra
type MessageBus interface {
	SendMessage(ctx context.Context, topicName string, message []byte) error
	CreateSubscription(ctx context.Context, topicName string) (*redis.PubSub, error)
}

type PubSub struct {
	MessageBus
}

func NewPubSub(db *redis.Client) *PubSub {
	return &PubSub{
		MessageBus: internal.NewMessageBus(db),
	}
}
