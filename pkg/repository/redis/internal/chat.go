package internal

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
)

type ChatCache struct {
	db *redis.Client
}

func NewChatCache(db *redis.Client) *ChatCache {
	return &ChatCache{db: db}
}

func (r *ChatCache) GetChatMembers(ctx context.Context, chatId int) ([]int, error) {
	log.Printf("Load members from cache")
	key := fmt.Sprintf("chat:%d:members", chatId)

	members, err := r.db.SMembers(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var memberIDs []int
	for _, member := range members {
		memberID, err := strconv.Atoi(member)
		if err != nil {
			return nil, err
		}
		memberIDs = append(memberIDs, memberID)
	}

	if len(members) != 0 {
		log.Printf("Load members from cache success")
	} else {
		log.Printf("Load members from cache failed")
	}
	return memberIDs, nil
}

func (r *ChatCache) SetChatMembers(ctx context.Context, chatId int, members []int) error {
	key := fmt.Sprintf("chat:%d:members", chatId)

	var memberStrings []string
	for _, memberID := range members {
		memberStrings = append(memberStrings, strconv.Itoa(memberID))
	}

	err := r.db.SAdd(ctx, key, memberStrings).Err()
	if err != nil {
		return err
	}

	return nil
}
