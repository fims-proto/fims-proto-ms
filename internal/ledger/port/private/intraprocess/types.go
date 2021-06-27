package intraprocess

import (
	commonaccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type VoucherLineItemRequest struct {
	AccountNumber string
	Debit         decimal.Decimal
	Credit        decimal.Decimal
}

type UpdateLedgerBalanceRequest struct {
	Sob         string
	VoucherUUID uuid.UUID
	LineItems   []VoucherLineItemRequest
}

type LoadLedgersRequest struct {
	Number         string
	Title          string
	SuperiorNumber string
	AccountType    commonaccount.Type
}
