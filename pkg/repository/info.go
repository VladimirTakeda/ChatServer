package repository

import (
	"ChatServer/pkg/types"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type InfoPostgres struct {
	db *pgxpool.Pool
}

func NewInfoPostgres(db *pgxpool.Pool) *InfoPostgres {
	return &InfoPostgres{db: db}
}

func (r *InfoPostgres) GetUsersByPrefix(ctx context.Context, prefix string) (types.UsersList, error) {
	//TODO add index for login
	createItemQuery := fmt.Sprintf("SELECT id, login FROM %s	WHERE login LIKE $1", userTable)

	rows, err := r.db.Query(ctx, createItemQuery, prefix+"%")
	if err != nil {
		logrus.Infof("row.Scan(&userId) %s", err.Error())
		return types.UsersList{}, err
	}

	var result types.UsersList

	for rows.Next() {
		userInfo := types.UserInfo{}
		err := rows.Scan(&userInfo.UserId, &userInfo.Nickname)
		if err != nil {
			return types.UsersList{}, err
		}
		result.Users = append(result.Users, userInfo)
	}

	if err := rows.Err(); err != nil {
		return types.UsersList{}, err
	}

	return result, nil
}
