package query

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Ledger struct {
	Id             uuid.UUID
	Sob            string
	Number         string
	Title          string
	SuperiorNumber string
	AccountType    string
	Debit          decimal.Decimal
	Credit         decimal.Decimal
	Balance        decimal.Decimal
}
