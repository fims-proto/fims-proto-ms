package db

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/counter/app/query"
	"github/fims-proto/fims-proto-ms/internal/counter/domain"

	"github.com/google/uuid"
)

type counter struct {
	Id             uuid.UUID `gorm:"type:uuid"`
	Current        uint
	Prefix         string
	Suffix         string
	BusinessObject string `gorm:"uniqueIndex"`
	LastResetAt    time.Time
	CreatedAt      time.Time `gorm:"<-:create"`
	UpdatedAt      time.Time
}

func marshall(c *domain.Counter) *counter {
	return &counter{
		Id:             c.Id(),
		Current:        c.CurrentIndex(),
		Prefix:         c.Prefix(),
		Suffix:         c.Suffix(),
		BusinessObject: c.BusinessObject(),
		LastResetAt:    c.LastResetAt(),
	}
}

func unmarshallToDomain(dbc *counter) (*domain.Counter, error) {
	return domain.NewCounterFromDB(dbc.Id, dbc.Current, dbc.BusinessObject, dbc.Prefix, dbc.Suffix, dbc.LastResetAt)
}

func unmarshallToQuery(dbc *counter) query.Counter {
	return query.Counter{
		Id: dbc.Id,
	}
}
