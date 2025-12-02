package query

import (
	"time"

	"github.com/google/uuid"
)

type Sob struct {
	Id                  uuid.UUID
	Name                string
	Description         string
	BaseCurrency        string
	StartingPeriodYear  int
	StartingPeriodMonth int
	AccountsCodeLength  []int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
