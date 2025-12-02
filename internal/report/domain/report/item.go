package report

import (
	"errors"

	commonerrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Item struct {
	id               uuid.UUID
	text             string
	level            int // starts from 1
	sequence         int // sequence within the parent, starts from 1
	sumFactor        int // 0,1,-1
	displaySumFactor bool
	dataSource       data_source.DataSource
	formulas         []*Formula
	amounts          []decimal.Decimal
	isEditable       bool
	isBreakdownItem  bool
	isAbleToAddChild bool
	isAbleToAddLeaf  bool
}

func NewItem(
	id uuid.UUID,
	text string,
	level int,
	sequence int,
	sumFactor int,
	displaySumFactor bool,
	dataSource string,
	formulas []*Formula,
	amounts []decimal.Decimal,
	isEditable bool,
	isBreakdownItem bool,
	isAbleToAddChild bool,
	isAbleToAddLeaf bool,
) (*Item, error) {
	if id == uuid.Nil {
		return nil, errors.New("item id cannot be nil")
	}

	if text == "" {
		return nil, commonerrors.NewSlugError("report-item-emptyText")
	}

	if level == 0 {
		return nil, commonerrors.NewSlugError("report-item-invalidLevel")
	}

	if sequence == 0 {
		return nil, commonerrors.NewSlugError("report-item-zeroSequence")
	}

	if sumFactor != -1 && sumFactor != 0 && sumFactor != 1 {
		return nil, commonerrors.NewSlugError("report-item-invalidSumFactor")
	}

	newDataSource, err := data_source.FromString(dataSource)
	if err != nil {
		return nil, err
	}

	if newDataSource != data_source.Formulas && len(formulas) > 0 {
		return nil, commonerrors.NewSlugError("report-item-invalidDataSourceWithFormulas")
	}

	if level == 1 && isBreakdownItem {
		return nil, commonerrors.NewSlugError("report-item-rootLevelIsBreakdownItem")
	}

	return &Item{
		id:               id,
		text:             text,
		level:            level,
		sequence:         sequence,
		sumFactor:        sumFactor,
		displaySumFactor: displaySumFactor,
		dataSource:       newDataSource,
		formulas:         formulas,
		amounts:          amounts,
		isEditable:       isEditable,
		isBreakdownItem:  isBreakdownItem,
		isAbleToAddChild: isAbleToAddChild,
		isAbleToAddLeaf:  isAbleToAddLeaf,
	}, nil
}

func (i *Item) copy() *Item {
	var newFormulas []*Formula
	for _, formula := range i.formulas {
		newFormulas = append(newFormulas, formula.copy())
	}

	newItem, _ := NewItem(
		uuid.New(),
		i.text,
		i.level,
		i.sequence,
		i.sumFactor,
		i.displaySumFactor,
		i.dataSource.String(),
		newFormulas,
		nil,
		i.isEditable,
		i.isBreakdownItem,
		i.isAbleToAddChild,
		i.isAbleToAddLeaf,
	)
	return newItem
}

func (i *Item) SetAmounts(amounts []decimal.Decimal) {
	i.amounts = amounts
}

func (i *Item) Id() uuid.UUID {
	return i.id
}

func (i *Item) Text() string {
	return i.text
}

func (i *Item) Level() int {
	return i.level
}

func (i *Item) Sequence() int {
	return i.sequence
}

func (i *Item) SumFactor() int {
	return i.sumFactor
}

func (i *Item) DisplaySumFactor() bool {
	return i.displaySumFactor
}

func (i *Item) DataSource() data_source.DataSource {
	return i.dataSource
}

func (i *Item) Formulas() []*Formula {
	return i.formulas
}

func (i *Item) Amounts() []decimal.Decimal {
	return i.amounts
}

func (i *Item) IsEditable() bool {
	return i.isEditable
}

func (i *Item) IsBreakdownItem() bool {
	return i.isBreakdownItem
}

func (i *Item) IsAbleToAddChild() bool {
	return i.isAbleToAddChild
}

func (i *Item) IsAbleToAddLeaf() bool {
	return i.isAbleToAddLeaf
}
