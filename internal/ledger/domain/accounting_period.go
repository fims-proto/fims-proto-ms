package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AccountingPeriod struct {
	id               uuid.UUID
	sobId            uuid.UUID
	previousPeriodId uuid.UUID
	financialYear    int
	number           int
	openingTime      time.Time
	endingTime       time.Time
	isClosed         bool
}

func NewAccountingPeriod(id, sobId, previousPeriodId uuid.UUID, financialYear, number int, openingTime, endingTime time.Time, isClosed bool) (*AccountingPeriod, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil accounting period id")
	}
	if sobId == uuid.Nil {
		return nil, errors.New("nil sob id")
	}
	if financialYear < 1970 || financialYear > 9999 {
		return nil, errors.New("invalid financial year")
	}
	if number < 1 || number > 12 {
		return nil, errors.New("invalid period number")
	}

	// TODO: question here: should ending time be provided when creating period?
	if openingTime.IsZero() {
		return nil, errors.New("zero opening time")
	}
	if !endingTime.IsZero() && openingTime.After(endingTime) {
		return nil, errors.New("opening time is after ending time")
	}

	return &AccountingPeriod{
		id:               id,
		sobId:            sobId,
		previousPeriodId: previousPeriodId,
		financialYear:    financialYear,
		number:           number,
		openingTime:      openingTime,
		endingTime:       endingTime,
		isClosed:         isClosed,
	}, nil
}

func (a AccountingPeriod) Id() uuid.UUID {
	return a.id
}

func (a AccountingPeriod) SobId() uuid.UUID {
	return a.sobId
}

func (a AccountingPeriod) PreviousPeriodId() uuid.UUID {
	return a.previousPeriodId
}

func (a AccountingPeriod) FinancialYear() int {
	return a.financialYear
}

func (a AccountingPeriod) Number() int {
	return a.number
}

func (a AccountingPeriod) OpeningTime() time.Time {
	return a.openingTime
}

func (a AccountingPeriod) EndingTime() time.Time {
	return a.endingTime
}

func (a AccountingPeriod) IsClosed() bool {
	return a.isClosed
}

func (a *AccountingPeriod) Close() error {
	if a.IsClosed() {
		return errors.New("period is already closed")
	}
	a.isClosed = true
	return nil
}

func (a *AccountingPeriod) Reopen() error {
	if !a.IsClosed() {
		return errors.New("period is not closed")
	}
	return errors.New("reopen period is not supported")
}
