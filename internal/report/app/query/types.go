package query

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// summary info when query to list all reports
type ReportInfo struct {
	Id         uuid.UUID
	TemplateId uuid.UUID
	PeriodId   uuid.UUID
	Title      string
	Level      int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Report struct {
	Id        uuid.UUID
	PeriodId  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Template  Template
}

type Account struct {
	Id               uuid.UUID
	SobId            uuid.UUID
	Title            string
	AccountNumber    string
	NumberHierarchy  []int
	Level            int
	Class            int
	Group            int
	BalanceDirection string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type TemplateInfo struct {
	Id        uuid.UUID
	SobId     uuid.UUID
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Template struct {
	Id        uuid.UUID
	SobId     uuid.UUID
	Title     string
	Tables    []Table
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Table struct {
	Header Header
	Items  []LineItem
}

type Header struct {
	Text    string
	Columns []string
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

type LineItem struct {
	Id                 uuid.UUID
	Text               string
	DataSource         string
	Formulas           []Formula
	Values             []decimal.Decimal
	SumFactor          int
	Level              int
	RowNumber          int
	Sequence           int
	DisplaySumFactor   bool
	DisplayRowNumber   bool
	IsDeletable        bool
	IsDraggable        bool
	IsAbleToAddChild   bool
	IsAbleToAddSibling bool
}

type Formula struct {
	Id               uuid.UUID
	AccountId        uuid.UUID
	ItemId           uuid.UUID
	isAccountFormula bool
	sumFactor        int
	rule             string
}

type User struct {
	Id     uuid.UUID
	Traits json.RawMessage
}
