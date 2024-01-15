package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

type DevicePostgres struct {
	db *pgxpool.Pool
}

func NewDevicePostgres(db *pgxpool.Pool) *DevicePostgres {
	return &DevicePostgres{db: db}
}

func (r *DevicePostgres) RegisterDevice(ctx context.Context, deviceId string, userId int) error {
	createItemQuery := fmt.Sprintf("INSERT INTO %s (deviceId, userId) values ($1, $2) RETURNING id", DevicesTable)

	row := r.db.QueryRow(ctx, createItemQuery, deviceId, userId)
	var id int
	err := row.Scan(&id)
	if err != nil {
		logrus.Infof("Failed to scan divece row:  %s", err.Error())
		return err
	}

	return nil
}

func (r *DevicePostgres) SaveLastActiveTime(ctx context.Context, userId int, deviceId string, lastTime time.Time) error {
	createItemQuery := fmt.Sprintf("UPDATE %s SET lastSeen = $1 WHERE deviceId = $2 AND userId = $3;", DevicesTable)

	rows, err := r.db.Query(ctx, createItemQuery, lastTime, deviceId, userId)
	if err != nil {
		logrus.Infof("Failed to execute query: %s", err.Error())
		return err
	}
	// It's important to keep available connection, rows.Close() 0 closes the current connection
	defer rows.Close()

	return nil
}

func (r *DevicePostgres) GetLastActiveTime(ctx context.Context, deviceId string, userId int) (time.Time, error) {
	createItemQuery := fmt.Sprintf("SELECT lastSeen FROM %s WHERE userId = $1 AND deviceId = $2", DevicesTable)

	row := r.db.QueryRow(ctx, createItemQuery, userId, deviceId)

	var lastSeen time.Time
	err := row.Scan(&lastSeen)
	if err != nil {
		logrus.Infof("row.Scan(lastSeen) %s", err.Error())
		return time.Time{}, err
	}

	return lastSeen, nil
}
