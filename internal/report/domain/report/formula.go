package report

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	commonerrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/formula_rule"
)

type Formula struct {
	id        uuid.UUID
	accountId uuid.UUID
	sumFactor int
	rule      formula_rule.FormulaRule
	amounts   []decimal.Decimal
}

func NewFormula(
	id uuid.UUID,
	accountId uuid.UUID,
	sumFactor int,
	rule string,
	amounts []decimal.Decimal,
) (*Formula, error) {
	if id == uuid.Nil {
		return nil, errors.New("formula id cannot be nil")
	}

	if accountId == uuid.Nil {
		return nil, commonerrors.NewSlugError("report-formula-emptyAccountId")
	}

	if sumFactor != 1 && sumFactor != -1 {
		return nil, commonerrors.NewSlugError("report-formula-invalidSumFactor")
	}

	newRule, err := formula_rule.FromString(rule)
	if err != nil {
		return nil, err
	}

	return &Formula{
		id:        id,
		accountId: accountId,
		sumFactor: sumFactor,
		rule:      newRule,
		amounts:   amounts,
	}, nil
}

func (f *Formula) copy() *Formula {
	newFormula, _ := NewFormula(
		uuid.New(),
		f.accountId,
		f.sumFactor,
		f.rule.String(),
		nil,
	)
	return newFormula
}

func (f *Formula) SetAmounts(amounts []decimal.Decimal) {
	f.amounts = amounts
}

func (f *Formula) Id() uuid.UUID {
	return f.id
}

func (f *Formula) AccountId() uuid.UUID {
	return f.accountId
}

func (f *Formula) SumFactor() int {
	return f.sumFactor
}

func (f *Formula) Rule() formula_rule.FormulaRule {
	return f.rule
}

func (f *Formula) Amounts() []decimal.Decimal {
	return f.amounts
}
