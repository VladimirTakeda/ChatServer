package service

import (
	"ChatServer/pkg/repository/postgres"
	"ChatServer/pkg/repository/redis"
	"ChatServer/pkg/types"
	"context"
	"time"
)

type ChatService struct {
	cache  redis.Chat
	chat   postgres.Chat
	device postgres.Device
}

func NewChatService(chatRepo postgres.Chat, deviceRepo postgres.Device, cacheRepo redis.Chat) *ChatService {
	return &ChatService{chat: chatRepo, device: deviceRepo, cache: cacheRepo}
}

func (s *ChatService) GetMissedMessages(ctx context.Context, userId int, deviceId string) ([]types.WsMessageOut, error) {
	// get all the chats
	chatIds, err := s.chat.GetAllChats(ctx, userId)
	if err != nil {
		return []types.WsMessageOut{}, err
	}

	var lastSeen time.Time
	lastSeen, err = s.device.GetLastActiveTime(ctx, deviceId, userId)
	if err != nil {
		return []types.WsMessageOut{}, err
	}

	// get last seen time
	return s.chat.GetMissedMessagesFromChats(ctx, userId, chatIds, lastSeen)
}

func (s *ChatService) CreateChat(ctx context.Context, users []int) (*int, error) {
	//TODO check if chat already exist
	chatId, err := s.chat.CreateChat(ctx, users)
	if err != nil {
		return chatId, err
	}
	err = s.cache.SetChatMembers(ctx, *chatId, users)
	if err != nil {
		return nil, err
	}
	return chatId, err
}

func (s *ChatService) DeleteChat(ctx context.Context, chatId int) error {
	return s.chat.DeleteChat(ctx, chatId)
}

func (s *ChatService) GetChatMembers(ctx context.Context, chatId int) ([]int, error) {
	members, err := s.cache.GetChatMembers(ctx, chatId)
	if err != nil {
		members, err = s.chat.GetChatMembers(ctx, chatId)
		if err == nil {
			err := s.cache.SetChatMembers(ctx, chatId, members)
			if err != nil {
				return nil, err
			}
		}
	}
	return members, err
}
