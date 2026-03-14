package command

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type JournalLineCmd struct {
	Id                 uuid.UUID
	Text               string
	AccountNumber      string
	Amount             decimal.Decimal
	DimensionOptionIds []uuid.UUID
}
