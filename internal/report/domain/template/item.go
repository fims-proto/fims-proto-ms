package report

import (
	"github/fims-proto/fims-proto-ms/internal/report/domain/template/data_source"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Item struct {
	id                    uuid.UUID
	parent_id             uuid.UUID // if parent_id == id it is a root item
	breakdonItems         []uuid.UUID
	text                  string
	dataSource            data_source.DataSource
	formulas              []Formula
	values                []Cell[decimal.Decimal]
	sumFactor             int
	level                 int
	rowNumber             int
	sequence              int
	displayBreakdownItems bool
	displaySumFactor      bool
	displayRowNumber      bool
	isDeletable           bool
	isDraggable           bool
	isAbleToAddChild      bool
	isAbleToAddSibling    bool
}
