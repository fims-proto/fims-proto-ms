package report

import (
	"errors"

	commonerrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/formula_rule"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Formula struct {
	id        uuid.UUID
	sequence  int // sequence within the parent, starts from 1
	accountId uuid.UUID
	sumFactor int
	rule      formula_rule.FormulaRule
	amounts   []decimal.Decimal
}

func NewFormula(
	id uuid.UUID,
	sequence int,
	accountId uuid.UUID,
	sumFactor int,
	rule string,
	amounts []decimal.Decimal,
) (*Formula, error) {
	if id == uuid.Nil {
		return nil, errors.New("formula id cannot be nil")
	}

	if sequence == 0 {
		return nil, commonerrors.NewSlugError("report-formula-zeroSequence")
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
		sequence:  sequence,
		accountId: accountId,
		sumFactor: sumFactor,
		rule:      newRule,
		amounts:   amounts,
	}, nil
}

func (f *Formula) copy() *Formula {
	newFormula, _ := NewFormula(
		uuid.New(),
		f.sequence,
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

func (f *Formula) Sequence() int {
	return f.sequence
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

// Equal checks if two Formula instances are equal by comparing all fields
func (f *Formula) Equal(other *Formula) bool {
	if f == nil || other == nil {
		return f == other
	}

	// Compare all fields except amounts (which is a slice and not directly comparable)
	if f.sequence != other.sequence ||
		f.accountId != other.accountId ||
		f.sumFactor != other.sumFactor ||
		f.rule != other.rule {
		return false
	}

	// Compare amounts slice
	if len(f.amounts) != len(other.amounts) {
		return false
	}
	for i := range f.amounts {
		if !f.amounts[i].Equal(other.amounts[i]) {
			return false
		}
	}

	return true
}
