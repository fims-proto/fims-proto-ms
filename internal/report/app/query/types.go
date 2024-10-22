package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Report struct {
	Id          uuid.UUID
	SobId       uuid.UUID
	Period      *Period
	Title       string
	Template    bool
	Class       string
	AmountTypes []string
	Sections    []Section

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Section struct {
	Id       uuid.UUID
	Title    string
	Amounts  []decimal.Decimal
	Sections []Section
	Items    []Item
}

type Item struct {
	Id               uuid.UUID
	Text             string
	Level            int
	SumFactor        int
	DisplaySumFactor bool
	DataSource       string
	Formulas         []Formula
	Amounts          []decimal.Decimal
	IsBreakdownItem  bool
	IsDeletable      bool
	IsTextModifiable bool
	IsDraggable      bool
	IsAbleToAddChild bool
	IsAbleToAddLeaf  bool
}

type Formula struct {
	Id        uuid.UUID
	Account   Account
	SumFactor int
	Rule      string
	Amounts   []decimal.Decimal
}

type Period struct {
	FiscalYear   int
	PeriodNumber int
}

type Account struct {
	Id                uuid.UUID
	SobId             uuid.UUID
	SuperiorAccountId *uuid.UUID
	Title             string
	AccountNumber     string
	Level             int
	IsLeaf            bool
	Class             int
	Group             int
	BalanceDirection  string
}
