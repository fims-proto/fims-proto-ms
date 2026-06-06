package query

import (
	"time"

	"github.com/google/uuid"
)

type DimensionCategory struct {
	Id        uuid.UUID
	SobId     uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DimensionOption struct {
	Id         uuid.UUID
	CategoryId uuid.UUID
	Name       string
	Category   DimensionCategory
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
