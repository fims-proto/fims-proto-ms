package datasource

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Connector is the basic component to get *gorm.DB by given DSN
type Connector struct{}

func NewConnector() Connector {
	return Connector{}
}

func (d Connector) GetConnection(dsn string) (*gorm.DB, error) {
	db, err := retry(100, 3*time.Second, func() (any, error) {
		return d.get(dsn)
	})
	if err != nil {
		return nil, err
	}
	return db.(*gorm.DB), nil
}

func (d Connector) get(dsn string) (*gorm.DB, error) {
	logLevel := logger.Warn
	if viper.GetBool("logger.showSql") {
		logLevel = logger.Info
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "a_",
		},
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(15 * time.Minute)
	return db, nil
}

func retry(tryTimes int, interval time.Duration, task func() (any, error)) (returning any, err error) {
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
