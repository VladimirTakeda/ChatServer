package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"time"
)

type ChatPostgres struct {
	db *pgxpool.Pool
}

func NewChatPostgres(db *pgxpool.Pool) *ChatPostgres {
	return &ChatPostgres{db: db}
}

func (r *ChatPostgres) CreateChat(ctx context.Context, users []int) (*int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		logrus.Infof("Failed to  %s", err.Error())
		return nil, err
	}

	var chatId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (users, createdTime) values ($1, $2) RETURNING id", chatTable)

	row := tx.QueryRow(ctx, createItemQuery, pq.Array(users), time.Now())
	err = row.Scan(&chatId)
	if err != nil {
		logrus.Infof("row.Scan(&chatId) %s", err.Error())
		tx.Rollback(ctx)
		return nil, err
	}

	createItemQuery = fmt.Sprintf("UPDATE %s SET chatIds = array_append(chatIds, $1) WHERE id = $2", userTable)

	for _, userId := range users {
		_, err = tx.Exec(ctx, createItemQuery, chatId, userId)
		if err != nil {
			logrus.Infof("Query failed: %s:", err.Error())
			tx.Rollback(ctx)
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		logrus.Infof("can't commit the transaction %s", err.Error())
		tx.Rollback(ctx)
		return nil, err
	}

	return &chatId, nil
}

func (r *ChatPostgres) DeleteChat(ctx context.Context, chatId int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		logrus.Infof("Failed to  %s", err.Error())
		return err
	}

	createItemQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $ RETURNING users;", chatTable)

	row := tx.QueryRow(ctx, createItemQuery, chatId)

	var users []int
	err = row.Scan(&users)
	if err != nil {
		logrus.Infof("row.Scan(&chatId) %s", err.Error())
		tx.Rollback(ctx)
		return err
	}

	createItemQuery = fmt.Sprintf("UPDATE %s SET ChatIds = array_remove(ChatIds, chatId) WHERE id = ANY($);", userTable)

	tx.Query(ctx, createItemQuery, users)
	if err != nil {
		logrus.Infof("row.Scan(&chatId) %s", err.Error())
		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logrus.Infof("can't commit the transaction %s", err.Error())
		tx.Rollback(ctx)
		return err
	}

	return nil
}
