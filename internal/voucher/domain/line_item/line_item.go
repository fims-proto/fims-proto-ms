package line_item

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
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
		return nil, commonErrors.NewSlugError("lineItem-emptyId", "empty item id")
	}

	if accountId == uuid.Nil {
		return nil, commonErrors.NewSlugError("lineItem-emptyAccountId", "empty account id")
	}

	if text == "" {
		return nil, commonErrors.NewSlugError("lineItem-emptyText", "empty line item text")
	}

	if debit.IsZero() && credit.IsZero() {
		return nil, commonErrors.NewSlugError("lineItem-emptyDebitCredit", "debit and credit are both zero")
	}

	if !debit.IsZero() && !credit.IsZero() {
		return nil, commonErrors.NewSlugError("lineItem-debitCreditDuplicated", "debit and credit are both presented")
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
