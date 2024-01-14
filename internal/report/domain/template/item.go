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

func NewLineItem(
	id uuid.UUID,
	text string,
	dataSource data_source.DataSource,
	formulas []Formula,
	values []decimal.Decimal,
	sumFactor int,
	level int,
	rowNumber int,
	sequence int,
	displaySumFactor bool,
	displayRowNumber bool,
	isDeletable bool,
	isDraggable bool,
	isAbleToAddChild bool,
	isAbleToAddSibling bool,
) (*Item, error) {
	var newValues []Cell[decimal.Decimal]
	for _, val := range values {
		cell := Cell[decimal.Decimal]{key: "", value: val}
		newValues = append(newValues, cell)
	}
	// Create a new Item instance
	item := &Item{
		id:                 id,
		text:               text,
		dataSource:         dataSource,
		formulas:           formulas,
		values:             newValues,
		sumFactor:          sumFactor,
		level:              level,
		rowNumber:          rowNumber,
		sequence:           sequence,
		displaySumFactor:   displaySumFactor,
		displayRowNumber:   displayRowNumber,
		isDeletable:        isDeletable,
		isDraggable:        isDraggable,
		isAbleToAddChild:   isAbleToAddChild,
		isAbleToAddSibling: isAbleToAddSibling,
	}

	// Return the new Item instance
	return item, nil
}
