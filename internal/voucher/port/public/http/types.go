package http

import (
	"time"
)

type sluggableErr interface {
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
	AccountNumber string
	Credit        string
	Debit         string
	Id            string
	Summary       string
}

type LineItemResponse struct {
	AccountNumber string
	Credit        string
	Debit         string
	Id            string
	Summary       string
}

type RecordVoucherRequest struct {
	AttachmentQuantity int
	Creator            string
	LineItems          []LineItemRequest
	VoucherType        string
}

type ReviewVoucherRequest struct {
	Reviewer string
}

type UpdateVoucherRequest struct {
	LineItems []LineItemRequest
}

type VoucherResponse struct {
	AttachmentQuantity int
	Auditor            string
	CreatedAt          time.Time
	Creator            string
	Credit             string
	Debit              string
	Id                 string
	IsAudited          bool
	IsPosted           bool
	IsReviewed         bool
	LineItems          []LineItemResponse
	Number             string
	Reviewer           string
	Sob                string
	Type               string
}
