package query

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/google/uuid"
)

type LineItem struct {
	Id        uuid.UUID
	AccountId uuid.UUID
	Summary   string
	Debit     decimal.Decimal
	Credit    decimal.Decimal
}

type Voucher struct {
	Id                 uuid.UUID
	SobId              uuid.UUID
	VoucherType        string
	Number             string
	AttachmentQuantity uint
	LineItems          []LineItem
	Debit              decimal.Decimal
	Credit             decimal.Decimal
	Creator            string
	Reviewer           string
	IsReviewed         bool
	Auditor            string
	IsAudited          bool
	IsPosted           bool
	TransactionTime    time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
