package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type LedgerLog struct {
	id              uuid.UUID
	voucherId       uuid.UUID
	postingId       uuid.UUID
	accountId       uuid.UUID
	transactionTime time.Time
	debit           decimal.Decimal
	credit          decimal.Decimal
}

func NewLedgerLog(id, postingId, accountId, voucherId uuid.UUID, transactionTime time.Time, debit, credit decimal.Decimal) (*LedgerLog, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil ledger id")
	}
	if postingId == uuid.Nil {
		return nil, errors.New("nil posting id")
	}
	if accountId == uuid.Nil {
		return nil, errors.New("nil account id")
	}
	if voucherId == uuid.Nil {
		return nil, errors.New("nil voucher id")
	}
	if transactionTime.IsZero() {
		return nil, errors.New("zero transaction time")
	}

	return &LedgerLog{
		id:              id,
		voucherId:       voucherId,
		postingId:       postingId,
		accountId:       accountId,
		transactionTime: transactionTime,
		debit:           debit,
		credit:          credit,
	}, nil
}

func (l LedgerLog) Id() uuid.UUID {
	return l.id
}

func (l LedgerLog) PostingId() uuid.UUID {
	return l.postingId
}

func (l LedgerLog) AccountId() uuid.UUID {
	return l.accountId
}

func (l LedgerLog) VoucherId() uuid.UUID {
	return l.voucherId
}

func (l LedgerLog) TransactionTime() time.Time {
	return l.transactionTime
}

func (l LedgerLog) Debit() decimal.Decimal {
	return l.debit
}

func (l LedgerLog) Credit() decimal.Decimal {
	return l.credit
}
