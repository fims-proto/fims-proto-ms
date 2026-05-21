package cash_flow_item

import (
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/cash_flow_item/category"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/cash_flow_item/direction"

	"github.com/google/uuid"
)

type CashFlowItem struct {
	id        uuid.UUID
	sobId     uuid.UUID
	code      string
	name      string
	category  category.Category
	direction direction.Direction
	sequence  int
}

func New(
	id uuid.UUID,
	sobId uuid.UUID,
	code string,
	name string,
	cat string,
	dir string,
	seq int,
) (*CashFlowItem, error) {
	c, err := category.FromString(cat)
	if err != nil {
		return nil, err
	}

	d, err := direction.FromString(dir)
	if err != nil {
		return nil, err
	}

	return &CashFlowItem{
		id:        id,
		sobId:     sobId,
		code:      code,
		name:      name,
		category:  c,
		direction: d,
		sequence:  seq,
	}, nil
}

func (i *CashFlowItem) Id() uuid.UUID {
	return i.id
}

func (i *CashFlowItem) SobId() uuid.UUID {
	return i.sobId
}

func (i *CashFlowItem) Code() string {
	return i.code
}

func (i *CashFlowItem) Name() string {
	return i.name
}

func (i *CashFlowItem) Category() category.Category {
	return i.category
}

func (i *CashFlowItem) Direction() direction.Direction {
	return i.direction
}

func (i *CashFlowItem) Sequence() int {
	return i.sequence
}
