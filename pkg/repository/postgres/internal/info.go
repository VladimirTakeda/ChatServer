package internal

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
	createItemQuery := fmt.Sprintf("SELECT id, login FROM %s	WHERE login LIKE $1", UserTable)

	rows, err := r.db.Query(ctx, createItemQuery, prefix+"%")
	if err != nil {
		logrus.Infof("row.Scan(&userId) %s", err.Error())
		return types.UsersList{}, err
	}
	defer rows.Close()

	var result types.UsersList

	for rows.Next() {
		userInfo := types.UserInfo{}
		err := rows.Scan(&userInfo.UserId, &userInfo.Nickname)
		if err != nil {
			logrus.Infof("Error scanning row: %s", err.Error())
			return types.UsersList{}, err
		}
		result.Users = append(result.Users, userInfo)
	}

	if err := rows.Err(); err != nil {
		logrus.Infof("Error reading rows: %s", err.Error())
		return types.UsersList{}, err
	}

	return result, nil
}
