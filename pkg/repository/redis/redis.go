package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	Host     string
	Port     string
	UserName string
	Password string
	DbNumber int
}

func NewRedisDb(cfg Config) (*redis.Client, error) {
	redisURL := fmt.Sprintf("redis://%s:%s@%s:%s/%d", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.DbNumber)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis connection config: %w", err)
	}

	rdb := redis.NewClient(opt)

	err = rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return rdb, nil
}
