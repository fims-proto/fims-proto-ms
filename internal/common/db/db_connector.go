package db

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DBConnector struct{}

func NewDBConnector() DBConnector {
	return DBConnector{}
}

func (d DBConnector) Open(username, password string) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		"localhost", 5432, username, password, "postgres")

	db, err := sqlx.Open("pgx", psqlInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open connection for user %s", username)
	}
	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}
