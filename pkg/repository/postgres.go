package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	chatTable             = "chats"
	userTable             = "users"
	couriersScheduleTable = "couriers_schedule"
)

type Config struct {
	Host     string
	Port     string
	UserName string
	Password string
	DbName   string
	SSLMode  string
}

func NewPostgresDb(cfg Config) (*pgxpool.Pool, error) {
	pgURL := fmt.Sprintf("postgres://postgres:%s@%s:%s/postgres?sslmode=disable", cfg.Password, cfg.Host, cfg.Port)

	var connErr error

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	connCfg, connErr := pgxpool.ParseConfig(pgURL)
	if connErr != nil {
		return nil, fmt.Errorf("failed to parse postgesql connection config: %w", connErr)
	}

	pgxPool, connErr := pgxpool.NewWithConfig(ctx, connCfg)
	if connErr != nil {
		return nil, connErr
	}

	connErr = pgxPool.Ping(ctx)
	if connErr != nil {
		return nil, connErr
	}
	return pgxPool, nil
}
