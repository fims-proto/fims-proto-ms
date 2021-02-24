package query

import (
	"time"
)

type LineItem struct {
	Summary       string
	AccountNumber string
	Debit         string
	Credit        string
}

type Voucher struct {
	UUID               string
	Number             string
	CreatedAt          time.Time
	AttachmentQuantity uint
	LineItems          []LineItem
	Debit              string
	Credit             string
	Creator            string
	Reviewer           string
	IsReviewed         bool
	Auditor            string
	IsAudited          bool
}
