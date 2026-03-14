package journal

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type JournalLine struct {
	id                 uuid.UUID
	accountId          uuid.UUID
	account            *account.Account
	text               string
	amount             decimal.Decimal
	dimensionOptionIds []uuid.UUID
}

func NewJournalLine(
	id uuid.UUID,
	account *account.Account,
	text string,
	amount decimal.Decimal,
	dimensionOptionIds []uuid.UUID,
) (*JournalLine, error) {
	if id == uuid.Nil {
		return nil, errors.NewSlugError("journalLine-emptyId")
	}

	if account == nil {
		return nil, errors.NewSlugError("journalLine-nilAccount")
	}

	if account.Id() == uuid.Nil {
		return nil, errors.NewSlugError("journalLine-emptyAccountId")
	}

	if text == "" {
		return nil, errors.NewSlugError("journalLine-emptyText")
	}

	if amount.IsZero() {
		return nil, errors.NewSlugError("journalLine-emptyAmount")
	}

	return &JournalLine{
		id:                 id,
		accountId:          account.Id(),
		account:            account,
		text:               text,
		amount:             amount,
		dimensionOptionIds: dimensionOptionIds,
	}, nil
}

func (i JournalLine) Id() uuid.UUID {
	return i.id
}

func (i JournalLine) AccountId() uuid.UUID {
	return i.accountId
}

func (i JournalLine) Account() *account.Account {
	return i.account
}

func (i JournalLine) Text() string {
	return i.text
}

func (i JournalLine) Amount() decimal.Decimal {
	return i.amount
}

func (i JournalLine) DimensionOptionIds() []uuid.UUID {
	return i.dimensionOptionIds
}
