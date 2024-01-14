package template

import (
	reportRule "github/fims-proto/fims-proto-ms/internal/report/domain/template/rule"

	"github.com/google/uuid"
)

type Formula struct {
	id               uuid.UUID
	accountId        uuid.UUID
	itemId           uuid.UUID
	isAccountFormula bool
	sumFactor        int
	rule             reportRule.Rule
}

func NewFormula(
	id uuid.UUID,
	accountId uuid.UUID,
	itemId uuid.UUID,
	isAccountFormula bool,
	sumFactor int,
	rule string,
) (*Formula, error) {
	newRule, err := reportRule.FromString(rule)
	if err != nil {
		return nil, err
	}

	return &Formula{
		id:               id,
		accountId:        accountId,
		itemId:           itemId,
		isAccountFormula: isAccountFormula,
		sumFactor:        sumFactor,
		rule:             newRule,
	}, nil
}

func (f *Formula) Id() uuid.UUID {
	return f.id
}

func (f *Formula) AccountId() uuid.UUID {
	return f.accountId
}

func (f *Formula) LineItemId() uuid.UUID {
	return f.itemId
}

func (f *Formula) IsAccountFormula() bool {
	return f.isAccountFormula
}

func (f *Formula) SumFactor() int {
	return f.sumFactor
}

func (f *Formula) Rule() reportRule.Rule {
	return f.rule
}
