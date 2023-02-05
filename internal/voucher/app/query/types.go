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
	Text          string
	Debit         decimal.Decimal
	Credit        decimal.Decimal
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type User struct {
	Id     uuid.UUID
	Traits json.RawMessage
}

type Period struct {
	PeriodId      uuid.UUID
	FinancialYear int
	Number        int
	OpeningTime   time.Time
	EndingTime    time.Time
	IsClosed      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Voucher struct {
	SobId              uuid.UUID
	Id                 uuid.UUID
	Period             Period
	VoucherType        string
	HeaderText         string
	DocumentNumber     string
	AttachmentQuantity int
	Debit              decimal.Decimal
	Credit             decimal.Decimal
	Creator            User
	Reviewer           User
	Auditor            User
	Poster             User
	IsReviewed         bool
	IsAudited          bool
	IsPosted           bool
	TransactionTime    time.Time
	LineItems          []LineItem
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
