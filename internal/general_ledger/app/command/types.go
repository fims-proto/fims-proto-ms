package command

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type JournalLineCmd struct {
	Id                uuid.UUID
	Text              string
	AccountNumber     string
	AuxiliaryAccounts []AuxiliaryItemCmd
	Amount            decimal.Decimal
}

type AuxiliaryItemCmd struct {
	CategoryKey string
	AccountKey  string
}
