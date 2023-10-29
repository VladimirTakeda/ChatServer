package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type UserPostgres struct {
	db *pgxpool.Pool
}

func NewUserPostgres(db *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Register(ctx context.Context, nickname string) (*int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		logrus.Infof("Failed to  %s", err.Error())
		return nil, err
	}

	var userId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (login) values ($1) RETURNING id", userTable)

	row := tx.QueryRow(ctx, createItemQuery, nickname)
	err = row.Scan(&userId)
	if err != nil {
		logrus.Infof("row.Scan(&userId) %s", err.Error())
		tx.Rollback(ctx)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logrus.Infof("can't commit the transaction %s", err.Error())
		tx.Rollback(ctx)
		return nil, err
	}

	return &userId, nil
}
