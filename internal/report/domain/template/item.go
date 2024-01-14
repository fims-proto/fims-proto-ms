package template

import (
	"github/fims-proto/fims-proto-ms/internal/report/domain/template/data_source"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Item struct {
	id                 uuid.UUID
	text               string
	dataSource         data_source.DataSource
	formulas           []Formula
	values             []Cell[decimal.Decimal]
	sumFactor          int
	level              int
	rowNumber          int
	sequence           int
	displaySumFactor   bool
	displayRowNumber   bool
	isDeletable        bool
	isDraggable        bool
	isAbleToAddChild   bool
	isAbleToAddSibling bool
}
