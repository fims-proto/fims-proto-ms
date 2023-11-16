package query

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"

	"github.com/google/uuid"
)

type Account struct {
	Id                  uuid.UUID
	SobId               uuid.UUID
	SuperiorAccountId   *uuid.UUID
	Title               string
	AccountNumber       string
	NumberHierarchy     []int
	Level               int
	AccountType         string
	BalanceDirection    string
	AuxiliaryCategories []AuxiliaryCategory
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type AuxiliaryCategory struct {
	Id         uuid.UUID
	SobId      uuid.UUID
	Key        string
	Title      string
	IsStandard bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type AuxiliaryAccount struct {
	Id          uuid.UUID
	Category    AuxiliaryCategory
	Key         string
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Period struct {
	Id           uuid.UUID
	SobId        uuid.UUID
	FiscalYear   int
	PeriodNumber int
	OpeningTime  time.Time
	EndingTime   time.Time
	IsClosed     bool
	IsCurrent    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Ledger struct {
	Id                   uuid.UUID
	SobId                uuid.UUID
	AccountId            uuid.UUID
	PeriodId             uuid.UUID
	OpeningDebitBalance  decimal.Decimal
	OpeningCreditBalance decimal.Decimal
	PeriodDebit          decimal.Decimal
	PeriodCredit         decimal.Decimal
	EndingDebitBalance   decimal.Decimal
	EndingCreditBalance  decimal.Decimal
	Account              Account
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type AuxiliaryLedger struct {
	Id                   uuid.UUID
	PeriodId             uuid.UUID
	AuxiliaryAccount     AuxiliaryAccount
	OpeningDebitBalance  decimal.Decimal
	OpeningCreditBalance decimal.Decimal
	PeriodDebit          decimal.Decimal
	PeriodCredit         decimal.Decimal
	EndingDebitBalance   decimal.Decimal
	EndingCreditBalance  decimal.Decimal
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type LineItem struct {
	Id                uuid.UUID
	Account           Account
	AuxiliaryAccounts []AuxiliaryAccount
	Text              string
	Debit             decimal.Decimal
	Credit            decimal.Decimal
	CreatedAt         time.Time
	UpdatedAt         time.Time
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
	Creator            *User
	Reviewer           *User
	Auditor            *User
	Poster             *User
	IsReviewed         bool
	IsAudited          bool
	IsPosted           bool
	TransactionTime    time.Time
	LineItems          []LineItem
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type User struct {
	Id     uuid.UUID
	Traits json.RawMessage
}
