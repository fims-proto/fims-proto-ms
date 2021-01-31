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
	CreatorUUID        string
	ReviewerUUID       string
	IsReviewed         bool
	AuditorUUID        string
	IsAudited          bool
}
