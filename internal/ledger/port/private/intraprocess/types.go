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
	VoucherUUID uuid.UUID
	LineItems   []VoucherLineItemRequest
}
