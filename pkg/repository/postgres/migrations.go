package postgres

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/sirupsen/logrus"
)

func RunMigrations(sourceUrl string, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logrus.Errorf("failed to get driver, error %s", err.Error())
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(sourceUrl, "postgres", driver)
	if err != nil {
		logrus.Errorf("failed to create migrations, error %s", err.Error())
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.Errorf("failed to UP migrations, error %s", err.Error())
		return err
	}
	return nil
}
