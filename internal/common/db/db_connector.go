package db

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnector struct{}

func NewDBConnector() DBConnector {
	return DBConnector{}
}

func (d DBConnector) Open(dsn string) (*gorm.DB, error) {
	db, err := retry(4, 5*time.Second, func() (interface{}, error) {
		return d.open(dsn)
	})
	if err != nil {
		return nil, err
	}
	return db.(*gorm.DB), nil
}

func (d DBConnector) open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open connection")
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

func retry(tryTimes int, interval time.Duration, task func() (interface{}, error)) (returning interface{}, err error) {
	for i := 0; i < tryTimes; i++ {
		returning, err = task()
		if err == nil {
			return returning, nil
		}
		if i < tryTimes-1 {
			time.Sleep(interval)
		}
	}
	return nil, err
}
