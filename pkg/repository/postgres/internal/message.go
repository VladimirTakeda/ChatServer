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

func (r *MessagePostgres) AddMessage(ctx context.Context, fromId, chatId int, text string, attachments []string) error {
	log.Printf("AddMessage from postgres")
	createItemQuery := fmt.Sprintf("INSERT INTO %s (createdTime, fromUserId, ChatId, Text, Attachments) values ($1, $2, $3, $4, $5)", MessageTable)

	_, err := r.db.Exec(ctx, createItemQuery, time.Now(), fromId, chatId, text, attachments)
	if err != nil {
		return err
	}

	log.Printf("AddMessage from postgres success")
	return nil
}
