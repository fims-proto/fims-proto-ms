package ledger

import (
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Ledger struct {
	id            uuid.UUID
	sobId         uuid.UUID
	periodId      uuid.UUID
	accountId     uuid.UUID
	account       *account.Account
	openingAmount decimal.Decimal
	periodAmount  decimal.Decimal
	periodDebit   decimal.Decimal
	periodCredit  decimal.Decimal
	endingAmount  decimal.Decimal
}

func New(
	id uuid.UUID,
	sobId uuid.UUID,
	periodId uuid.UUID,
	accountId uuid.UUID,
	account *account.Account,
	openingAmount decimal.Decimal,
	periodAmount decimal.Decimal,
	periodDebit decimal.Decimal,
	periodCredit decimal.Decimal,
	endingAmount decimal.Decimal,
) (*Ledger, error) {
	if id == uuid.Nil {
		return nil, commonErrors.NewInternalError(commonErrors.SlugLedgerNilId)
	}

	if sobId == uuid.Nil {
		return nil, commonErrors.NewInternalError(commonErrors.SlugLedgerNilSobId)
	}

	if periodId == uuid.Nil {
		return nil, commonErrors.NewInternalError(commonErrors.SlugLedgerNilPeriodId)
	}

	if accountId == uuid.Nil {
		return nil, commonErrors.NewInternalError(commonErrors.SlugLedgerNilAccountId)
	}

	if account == nil {
		return nil, commonErrors.NewInternalError(commonErrors.SlugLedgerNilAccount)
	}

	return &Ledger{
		id:            id,
		sobId:         sobId,
		accountId:     accountId,
		periodId:      periodId,
		openingAmount: openingAmount,
		periodAmount:  periodAmount,
		periodDebit:   periodDebit,
		periodCredit:  periodCredit,
		endingAmount:  endingAmount,
		account:       account,
	}, nil
}

func (l *Ledger) Id() uuid.UUID {
	return l.id
}

func (l *Ledger) SobId() uuid.UUID {
	return l.sobId
}

func (l *Ledger) PeriodId() uuid.UUID {
	return l.periodId
}

func (l *Ledger) AccountId() uuid.UUID {
	return l.accountId
}

func (l *Ledger) Account() *account.Account {
	return l.account
}

func (l *Ledger) OpeningAmount() decimal.Decimal {
	return l.openingAmount
}

func (l *Ledger) PeriodAmount() decimal.Decimal {
	return l.periodAmount
}

func (l *Ledger) PeriodDebit() decimal.Decimal {
	return l.periodDebit
}

func (l *Ledger) PeriodCredit() decimal.Decimal {
	return l.periodCredit
}

func (l *Ledger) EndingAmount() decimal.Decimal {
	return l.endingAmount
}
