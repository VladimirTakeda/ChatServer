package internal

import (
	"ChatServer/pkg/types"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
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
	createItemQuery := fmt.Sprintf("INSERT INTO %s (users, createdTime) values ($1, $2) RETURNING id", ChatTable)

	row := tx.QueryRow(ctx, createItemQuery, pq.Array(users), time.Now())
	err = row.Scan(&chatId)
	if err != nil {
		logrus.Infof("row.Scan(&chatId) %s", err.Error())
		tx.Rollback(ctx)
		return nil, err
	}

	createItemQuery = fmt.Sprintf("UPDATE %s SET chatIds = array_append(chatIds, $1) WHERE id = $2", UserTable)

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

	createItemQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $ RETURNING users;", ChatTable)

	row := tx.QueryRow(ctx, createItemQuery, chatId)

	var users []int
	err = row.Scan(&users)
	if err != nil {
		logrus.Infof("row.Scan(&chatId) %s", err.Error())
		tx.Rollback(ctx)
		return err
	}

	createItemQuery = fmt.Sprintf("UPDATE %s SET ChatIds = array_remove(ChatIds, chatId) WHERE id = ANY($);", UserTable)

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

func (r *ChatPostgres) GetAllChats(ctx context.Context, userId int) ([]int, error) {
	createItemQuery := fmt.Sprintf("SELECT ChatIds FROM %s WHERE id = $1", UserTable)

	row := r.db.QueryRow(ctx, createItemQuery, userId)

	var users []int
	err := row.Scan(&users)
	if err != nil {
		logrus.Infof("row.Scan(&ChatIds) %s", err.Error())
		return []int{}, err
	}

	return users, nil
}

func (r *ChatPostgres) GetMissedMessagesFromChats(ctx context.Context, userTo int, chatIds []int, lastSeen time.Time) ([]types.WsMessageOut, error) {
	// SQL запрос с использованием IN оператора и сортировкой
	query := fmt.Sprintf("SELECT createdTime, fromUserId, text, ChatId, Attachments FROM %s WHERE ChatId = ANY($1) AND createdTime > $2 ORDER BY ChatId, createdTime",
		MessageTable)

	rows, err := r.db.Query(ctx, query, chatIds, lastSeen)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []types.WsMessageOut

	for rows.Next() {
		messageInfo := types.WsMessageOut{}
		err := rows.Scan(&messageInfo.Time, &messageInfo.UserFrom, &messageInfo.Content, &messageInfo.ChatTo, &messageInfo.Attachments)
		if err != nil {
			return nil, err
		}

		result = append(result, messageInfo)
	}

	return result, nil
}

func (r *ChatPostgres) GetChatMembers(ctx context.Context, chatId int) ([]int, error) {
	log.Printf("Load members from postgres")
	createItemQuery := fmt.Sprintf("SELECT Users FROM %s WHERE id = $1", ChatTable)

	row := r.db.QueryRow(ctx, createItemQuery, chatId)

	var members []int
	err := row.Scan(&members)
	if err != nil {
		logrus.Infof("row.Scan(&members) %s", err.Error())
		return []int{}, err
	}

	log.Printf("Load members from postgres success")
	return members, nil
}
