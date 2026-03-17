package query

import (
	"encoding/json"
	"time"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github.com/shopspring/decimal"

	"github.com/google/uuid"
)

type DimensionCategory struct {
	Id   uuid.UUID
	Name string
}

type DimensionOption struct {
	Id         uuid.UUID
	Name       string
	CategoryId uuid.UUID
	Category   DimensionCategory
}

type Account struct {
	Id                   uuid.UUID
	SobId                uuid.UUID
	SuperiorAccountId    *uuid.UUID
	Title                string
	AccountNumber        string
	NumberHierarchy      []int
	Level                int
	IsLeaf               bool
	Class                int
	Group                int
	BalanceDirection     string
	DimensionCategoryIds []uuid.UUID         // internal: used by enricher, not exposed in HTTP response
	DimensionCategories  []DimensionCategory // populated by enricher on detail queries only
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type Period struct {
	Id           uuid.UUID
	SobId        uuid.UUID
	FiscalYear   int
	PeriodNumber int
	IsClosed     bool
	IsCurrent    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Ledger struct {
	Id            uuid.UUID
	SobId         uuid.UUID
	AccountId     uuid.UUID
	PeriodId      uuid.UUID
	OpeningAmount decimal.Decimal
	PeriodAmount  decimal.Decimal
	PeriodDebit   decimal.Decimal
	PeriodCredit  decimal.Decimal
	EndingAmount  decimal.Decimal
	Account       Account
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type JournalLine struct {
	Id                 uuid.UUID
	Account            Account
	Text               string
	Amount             decimal.Decimal
	DimensionOptionIds []uuid.UUID       // internal: used by enricher, not exposed in HTTP response
	DimensionOptions   []DimensionOption // populated by enricher on detail queries only
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type Journal struct {
	SobId              uuid.UUID
	Id                 uuid.UUID
	Period             Period
	HeaderText         string
	DocumentNumber     string
	JournalType        string
	ReferenceJournalId *uuid.UUID
	AttachmentQuantity int
	Amount             decimal.Decimal
	Creator            *User
	Reviewer           *User
	Auditor            *User
	Poster             *User
	IsReviewed         bool
	IsAudited          bool
	IsPosted           bool
	TransactionDate    transaction_date.TransactionDate
	JournalLines       []JournalLine
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type User struct {
	Id     uuid.UUID
	Traits json.RawMessage
}

type LedgerEntry struct {
	JournalId       uuid.UUID
	JournalNumber   string
	TransactionDate transaction_date.TransactionDate
	Text            string
	Amount          decimal.Decimal
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type LedgerDimensionSummaryItem struct {
	DimensionOptionId   uuid.UUID
	DimensionOptionName string
	TotalAmount         decimal.Decimal
}
