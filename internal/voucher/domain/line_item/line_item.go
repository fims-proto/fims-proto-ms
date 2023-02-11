package line_item

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

type LineItem struct {
	id        uuid.UUID
	accountId uuid.UUID
	text      string
	debit     decimal.Decimal
	credit    decimal.Decimal
}

func New(itemId, accountId uuid.UUID, text string, debit, credit decimal.Decimal) (*LineItem, error) {
	if itemId == uuid.Nil {
		return nil, errors.NewSlugError("lineItem-emptyId")
	}

	if accountId == uuid.Nil {
		return nil, errors.NewSlugError("lineItem-emptyAccountId")
	}

	if text == "" {
		return nil, errors.NewSlugError("lineItem-emptyText")
	}

	if debit.IsZero() && credit.IsZero() {
		return nil, errors.NewSlugError("lineItem-emptyDebitCredit")
	}

	if !debit.IsZero() && !credit.IsZero() {
		return nil, errors.NewSlugError("lineItem-debitCreditDuplicated")
	}

	return &LineItem{
		id:        itemId,
		accountId: accountId,
		text:      text,
		debit:     debit,
		credit:    credit,
	}, nil
}

func (i LineItem) Id() uuid.UUID {
	return i.id
}

func (i LineItem) AccountId() uuid.UUID {
	return i.accountId
}

func (i LineItem) Text() string {
	return i.text
}

func (i LineItem) Debit() decimal.Decimal {
	return i.debit
}

func (i LineItem) Credit() decimal.Decimal {
	return i.credit
}
