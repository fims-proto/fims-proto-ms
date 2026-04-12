package report

import (
	commonerrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/item_type"

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
	itemType         item_type.ItemType
	dataSource       data_source.DataSource
	formulas         []*Formula
	amounts          []decimal.Decimal
	isEditable       bool
	isBreakdownItem  bool
	isAbleToAddChild bool
}

func NewItem(
	id uuid.UUID,
	text string,
	level int,
	sequence int,
	itemType string,
	sumFactor int,
	displaySumFactor bool,
	dataSource string,
	formulas []*Formula,
	amounts []decimal.Decimal,
	isEditable bool,
	isBreakdownItem bool,
	isAbleToAddChild bool,
) (*Item, error) {
	if id == uuid.Nil {
		return nil, commonerrors.NewInternalError(commonerrors.SlugReportItemNilId)
	}

	if text == "" {
		return nil, commonerrors.NewInvalidInputError(commonerrors.SlugReportItemEmptyText)
	}

	if level == 0 {
		return nil, commonerrors.NewInvalidInputError(commonerrors.SlugReportItemInvalidLevel)
	}

	if sequence == 0 {
		return nil, commonerrors.NewInvalidInputError(commonerrors.SlugReportItemZeroSeq)
	}

	if sumFactor != -1 && sumFactor != 0 && sumFactor != 1 {
		return nil, commonerrors.NewInvalidInputError(commonerrors.SlugReportItemInvalidSumFactor)
	}

	newItemType, err := item_type.FromString(itemType)
	if err != nil {
		return nil, err
	}

	newDataSource, err := data_source.FromString(dataSource)
	if err != nil {
		return nil, err
	}

	if newDataSource != data_source.Formulas && len(formulas) > 0 {
		return nil, commonerrors.NewInvalidInputError(commonerrors.SlugReportItemInvalidDataSourceWithForms)
	}

	if level == 1 && isBreakdownItem {
		return nil, commonerrors.NewInvalidInputError(commonerrors.SlugReportItemRootIsBreakdown)
	}

	if isBreakdownItem && isAbleToAddChild {
		return nil, commonerrors.NewInvalidInputError(commonerrors.SlugReportItemBreakdownNoChild)
	}

	return &Item{
		id:               id,
		text:             text,
		level:            level,
		sequence:         sequence,
		sumFactor:        sumFactor,
		displaySumFactor: displaySumFactor,
		itemType:         newItemType,
		dataSource:       newDataSource,
		formulas:         formulas,
		amounts:          amounts,
		isEditable:       isEditable,
		isBreakdownItem:  isBreakdownItem,
		isAbleToAddChild: isAbleToAddChild,
	}, nil
}

func (i *Item) copy() *Item {
	var newFormulas []*Formula
	for _, formula := range i.formulas {
		newFormulas = append(newFormulas, formula.copy())
	}

	newItem, _ := NewItem(uuid.New(), i.text, i.level, i.sequence, i.itemType.String(), i.sumFactor, i.displaySumFactor, i.dataSource.String(), newFormulas, nil, i.isEditable, i.isBreakdownItem, i.isAbleToAddChild)
	return newItem
}

func (i *Item) SetAmounts(amounts []decimal.Decimal) {
	i.amounts = amounts
}

func (i *Item) SetSequence(sequence int) {
	i.sequence = sequence
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

func (i *Item) ItemType() item_type.ItemType {
	return i.itemType
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
