package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type MessagePostgres struct {
	db *pgxpool.Pool
}

func NewMessagePostgres(db *pgxpool.Pool) *MessagePostgres {
	return &MessagePostgres{db: db}
}

func (r *MessagePostgres) AddMessage(ctx context.Context, fromId, chatId int, text string) error {
	log.Printf("AddMessage from postgres")
	createItemQuery := fmt.Sprintf("INSERT INTO %s (createdTime, fromUserId, ChatId, Text) values ($1, $2, $3, $4)", MessageTable)

	_, err := r.db.Exec(ctx, createItemQuery, time.Now(), fromId, chatId, text)
	if err != nil {
		return err
	}

	log.Printf("AddMessage from postgres success")
	return nil
}
