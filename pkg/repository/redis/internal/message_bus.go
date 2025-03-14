package internal

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type MessageBusRedis struct {
	db *redis.Client
}

// TODO: replace pubsub to some interface
func (r *MessageBusRedis) CreateSubscription(ctx context.Context, topicName string) (*redis.PubSub, error) {
	return r.db.Subscribe(ctx, topicName), nil
}

func NewMessageBus(db *redis.Client) *MessageBusRedis {
	return &MessageBusRedis{db: db}
}

func (r *MessageBusRedis) SendMessage(ctx context.Context, topicName string, message []byte) error {
	result, err := r.db.Publish(ctx, topicName, message).Result()
	if err != nil {
		return err
	}
	fmt.Printf("Published message '%s' to topic '%s', received by %d subscribers\n", message, topicName, result)
	return nil
}
