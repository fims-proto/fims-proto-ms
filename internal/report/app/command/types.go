package command

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TemplateCmd struct {
	Id     uuid.UUID
	SobId  uuid.UUID
	Title  string
	Tables []TableCmd
}

type TableCmd struct {
	Header HeaderCmd
	Items  []LineItemCmd
}

type HeaderCmd struct {
	Text    string
	Columns []string
}

type LineItemCmd struct {
	Id                 uuid.UUID
	Text               string
	DataSource         string
	Formulas           []FormulaCmd
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

type FormulaCmd struct {
	AccountId        uuid.UUID
	ItemId           uuid.UUID
	isAccountFormula bool
	sumFactor        int
	rule             string
}
