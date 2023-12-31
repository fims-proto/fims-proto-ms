package report

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
)

type Item struct {
	id                 uuid.UUID
	text               string
	values             []Cell[decimal.Decimal]
	level              int
	sumFactor          int
	rowNumber          int
	sequence           int
	dataSource         data_source.DataSource
	formulas           []Formula
	displaySumFactor   bool
	displayRowNumber   bool
	isBreakdownItem    bool
	isDeletable        bool
	isNameModifiable   bool
	isDraggable        bool
	isAbleToAddChild   bool
	isAbleToAddSibling bool
	isAbleToAddLeaf    bool
}
