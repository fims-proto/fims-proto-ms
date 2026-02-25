package ledger_entry

import (
	"errors"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LedgerEntry struct {
	id                uuid.UUID
	sobId             uuid.UUID
	periodId          uuid.UUID
	voucherId         uuid.UUID
	lineItemId        uuid.UUID
	accountId         uuid.UUID
	auxiliaryAccounts []*auxiliary_account.AuxiliaryAccount
	transactionDate   transaction_date.TransactionDate
	amount            decimal.Decimal
}

func New(
	id uuid.UUID,
	sobId uuid.UUID,
	periodId uuid.UUID,
	voucherId uuid.UUID,
	lineItemId uuid.UUID,
	accountId uuid.UUID,
	auxiliaryAccounts []*auxiliary_account.AuxiliaryAccount,
	transactionDate transaction_date.TransactionDate,
	amount decimal.Decimal,
) (*LedgerEntry, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil ledger entry id")
	}

	if sobId == uuid.Nil {
		return nil, errors.New("nil sob id")
	}

	if periodId == uuid.Nil {
		return nil, errors.New("nil period id")
	}

	if voucherId == uuid.Nil {
		return nil, errors.New("nil voucher id")
	}

	if lineItemId == uuid.Nil {
		return nil, errors.New("nil line item id")
	}

	if accountId == uuid.Nil {
		return nil, errors.New("nil account id")
	}

	if transactionDate.IsZero() {
		return nil, commonErrors.NewSlugError("ledgerEntry-zeroTransactionDate")
	}

	if amount.IsZero() {
		return nil, commonErrors.NewSlugError("ledgerEntry-emptyAmount")
	}

	return &LedgerEntry{
		id:                id,
		sobId:             sobId,
		periodId:          periodId,
		voucherId:         voucherId,
		lineItemId:        lineItemId,
		accountId:         accountId,
		auxiliaryAccounts: auxiliaryAccounts,
		transactionDate:   transactionDate,
		amount:            amount,
	}, nil
}

func (l *LedgerEntry) Id() uuid.UUID {
	return l.id
}

func (l *LedgerEntry) SobId() uuid.UUID {
	return l.sobId
}

func (l *LedgerEntry) PeriodId() uuid.UUID {
	return l.periodId
}

func (l *LedgerEntry) VoucherId() uuid.UUID {
	return l.voucherId
}

func (l *LedgerEntry) LineItemId() uuid.UUID {
	return l.lineItemId
}

func (l *LedgerEntry) AccountId() uuid.UUID {
	return l.accountId
}

func (l *LedgerEntry) AuxiliaryAccounts() []*auxiliary_account.AuxiliaryAccount {
	return l.auxiliaryAccounts
}

func (l *LedgerEntry) TransactionDate() transaction_date.TransactionDate {
	return l.transactionDate
}

func (l *LedgerEntry) Amount() decimal.Decimal {
	return l.amount
}
