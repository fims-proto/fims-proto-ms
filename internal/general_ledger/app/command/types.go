package command

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LineItemCmd struct {
	Id                   uuid.UUID
	Text                 string
	AccountNumber        string
	AuxiliaryAccountKeys []string
	Debit                decimal.Decimal
	Credit               decimal.Decimal
}
