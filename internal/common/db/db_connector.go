package db

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnector struct{}

func NewDBConnector() DBConnector {
	return DBConnector{}
}

func (d DBConnector) Open(username, password string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s TimeZone=%s",
		"localhost", 5432, username, password, "postgres", "Asia/Shanghai")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open connection for user %s", username)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sql.DB")
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(15 * time.Minute)
	return db, nil
}
