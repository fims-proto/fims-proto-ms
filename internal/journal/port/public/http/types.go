package http

import (
	"time"

	"github.com/shopspring/decimal"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type AuditVoucherRequest struct {
	Auditor string `json:"auditor"`
}

type LineItemRequest struct {
	Id            string          `json:"id"`
	AccountNumber string          `json:"accountNumber"`
	Summary       string          `json:"summary"`
	Credit        decimal.Decimal `json:"credit"`
	Debit         decimal.Decimal `json:"debit"`
}

type LineItemResponse struct {
	Id            string          `json:"id"`
	AccountId     string          `json:"accountId"`
	AccountNumber string          `json:"accountNumber"`
	Summary       string          `json:"summary"`
	Credit        decimal.Decimal `json:"credit"`
	Debit         decimal.Decimal `json:"debit"`
}

type CreateVoucherRequest struct {
	AttachmentQuantity int               `json:"attachmentQuantity"`
	Creator            string            `json:"creator"`
	VoucherType        string            `json:"voucherType"`
	TransactionTime    time.Time         `json:"transactionTime"`
	LineItems          []LineItemRequest `json:"lineItems"`
}

type ReviewVoucherRequest struct {
	Reviewer string `json:"reviewer"`
}

type UpdateVoucherRequest struct {
	TransactionTime time.Time         `json:"transactionTime"`
	LineItems       []LineItemRequest `json:"lineItems"`
}

type UserResponse struct {
	Id     string `json:"id"`
	Traits any    `json:"traits"`
}

type PeriodResponse struct {
	Id            string    `json:"id"`
	FinancialYear int       `json:"financialYear"`
	Number        int       `json:"number"`
	OpeningTime   time.Time `json:"openingTime"`
	EndingTime    time.Time `json:"endingTime"`
	IsClosed      bool      `json:"isClosed"`
}

type VoucherResponse struct {
	Id                 string `json:"id"`
	SobId              string `json:"sobId"`
	Period             PeriodResponse
	Number             string             `json:"number"`
	Type               string             `json:"type"`
	AttachmentQuantity int                `json:"attachmentQuantity"`
	Auditor            UserResponse       `json:"auditor"`
	Creator            UserResponse       `json:"creator"`
	Reviewer           UserResponse       `json:"reviewer"`
	Credit             decimal.Decimal    `json:"credit"`
	Debit              decimal.Decimal    `json:"debit"`
	IsAudited          bool               `json:"isAudited"`
	IsPosted           bool               `json:"isPosted"`
	IsReviewed         bool               `json:"isReviewed"`
	TransactionTime    time.Time          `json:"transactionTime"`
	LineItems          []LineItemResponse `json:"lineItems"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}
