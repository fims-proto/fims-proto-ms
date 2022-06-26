package query

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"

	"github.com/google/uuid"
)

type LineItem struct {
	Id            uuid.UUID
	AccountId     uuid.UUID
	AccountNumber string
	Summary       string
	Debit         decimal.Decimal
	Credit        decimal.Decimal
}

type User struct {
	Id     uuid.UUID
	Traits json.RawMessage
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
	Creator            User
	Reviewer           User
	IsReviewed         bool
	Auditor            User
	IsAudited          bool
	IsPosted           bool
	TransactionTime    time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
