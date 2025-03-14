package service

import (
	"ChatServer/pkg/repository/postgres"
	"ChatServer/pkg/repository/redis"
	"ChatServer/pkg/types"
	"context"
	"encoding/json"
	redis2 "github.com/redis/go-redis/v9"
	"strconv"
)

type MessageService struct {
	repo       postgres.Message
	messageBus redis.MessageBus
}

func NewMessageService(repo postgres.Message, messageBus redis.MessageBus) *MessageService {
	return &MessageService{repo: repo, messageBus: messageBus}
}

// The function saves message in DB and send it to Redis PubSub
func (s *MessageService) AddMessage(ctx context.Context, message types.WsMessageWithTime, userIDs []int) error {
	err := s.repo.AddMessage(ctx, message.UserFrom, message.ChatTo, message.Content, message.Attachments)
	if err != nil {
		return err
	}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for _, userID := range userIDs {
		err = s.messageBus.SendMessage(ctx, "user:"+strconv.Itoa(userID), messageBytes)
	}
	return err
}

func (s *MessageService) CreateSubscription(ctx context.Context, topicName string) (*redis2.PubSub, error) {
	return s.messageBus.CreateSubscription(ctx, topicName)
}
