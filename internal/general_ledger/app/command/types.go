package command

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LineItemCmd struct {
	Id                uuid.UUID
	Text              string
	AccountNumber     string
	AuxiliaryAccounts []AuxiliaryItemCmd
	Debit             decimal.Decimal
	Credit            decimal.Decimal
}

type AuxiliaryItemCmd struct {
	CategoryKey string
	AccountKey  string
}
