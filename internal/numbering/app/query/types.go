package query

import (
	"time"

	"github.com/google/uuid"
)

type Identifier struct {
	Id                        uuid.UUID
	IdentifierConfigurationId uuid.UUID
	Identifier                string
	CreatedAt                 time.Time
}

type PropertyMatcher struct {
	Name  string
	Value string
}

type IdentifierConfiguration struct {
	Id                   uuid.UUID
	TargetBusinessObject string
	PropertyMatchers     []PropertyMatcher
	Counter              int
	Prefix               string
	Suffix               string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
