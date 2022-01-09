package http

import (
	"time"

	"github.com/shopspring/decimal"
)

type slugErr interface {
	Slug() string
}

type Error struct {
	Message string
	Slug    string
}

type AuditVoucherRequest struct {
	Auditor string
}

type LineItemRequest struct {
	Id            string
	AccountNumber string
	Summary       string
	Credit        string
	Debit         string
}

type LineItemResponse struct {
	Id        string
	AccountId string
	Summary   string
	Credit    decimal.Decimal
	Debit     decimal.Decimal
}

type CreateVoucherRequest struct {
	AttachmentQuantity int
	Creator            string
	VoucherType        string
	TransactionTime    time.Time
	LineItems          []LineItemRequest
}

type ReviewVoucherRequest struct {
	Reviewer string
}

type UpdateVoucherRequest struct {
	TransactionTime time.Time
	LineItems       []LineItemRequest
}

type VoucherResponse struct {
	Id                 string
	SobId              string
	Number             string
	Type               string
	AttachmentQuantity int
	Auditor            string
	Creator            string
	Credit             decimal.Decimal
	Debit              decimal.Decimal
	IsAudited          bool
	IsPosted           bool
	IsReviewed         bool
	Reviewer           string
	TransactionTime    time.Time
	LineItems          []LineItemResponse
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
