package query

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Report struct {
	Id       uuid.UUID
	SobId    uuid.UUID
	Period   *Period
	Title    string
	Template bool
	Class    string
	Columns  []Column
	Rows     []Row

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Column struct {
	Id        uuid.UUID
	Label     string
	ValueType string
	Sequence  int
}

type Row struct {
	Id          uuid.UUID
	RowCode     string
	Text        string
	Sequence    int
	LineNo      *int
	ShowLineNo  bool
	SumFactor   int
	CanEdit     bool
	CanMove     bool
	CanAddChild bool
	Expression  Expression
	Rows        []Row
	Amounts     []decimal.Decimal
}

type Expression struct {
	Id             uuid.UUID
	Kind           string
	LedgerAccounts []LedgerAccountReference
	CashFlowItems  []CashFlowItemReference
	RowReferences  []RowReference
}

type LedgerAccountReference struct {
	RawAccountNumber string
	AccountId        uuid.UUID
	SumFactor        int
	Measure          string
}

type CashFlowItemReference struct {
	Code      string
	ItemId    uuid.UUID
	SumFactor int
}

type RowReference struct {
	RowCode   string
	SumFactor int
}

type Period struct {
	FiscalYear   int
	PeriodNumber int
}
