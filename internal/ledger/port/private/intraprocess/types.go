package intraprocess

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type VoucherLineItemRequest struct {
	AccountNumber string
	Debit         decimal.Decimal
	Credit        decimal.Decimal
}

type UpdateLedgerBalanceRequest struct {
	Sob       string
	VoucherId uuid.UUID
	LineItems []VoucherLineItemRequest
}

type LoadLedgersRequest struct {
	Number         string
	Title          string
	SuperiorNumber string
	AccountType    string
}
