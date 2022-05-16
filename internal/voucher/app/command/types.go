package command

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LineItemCmd struct {
	Id            uuid.UUID
	Summary       string
	AccountNumber string
	Debit         decimal.Decimal
	Credit        decimal.Decimal
}
