package period

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Period struct {
	sobId            uuid.UUID
	periodId         uuid.UUID
	previousPeriodId uuid.UUID
	financialYear    int
	number           int
	openingTime      time.Time
	endingTime       time.Time
	isClosed         bool
}

func New(sobId, periodId, previousPeriodId uuid.UUID, financialYear, number int, openingTime, endingTime time.Time, isClosed bool) (*Period, error) {
	if sobId == uuid.Nil {
		return nil, errors.New("nil sob id")
	}

	if periodId == uuid.Nil {
		return nil, errors.New("nil period id")
	}

	if financialYear < 1970 || financialYear > 9999 {
		return nil, errors.Errorf("invalid financial year %d", financialYear)
	}

	if number < 1 {
		return nil, errors.Errorf("invalid period number %d", number)
	}

	if openingTime.IsZero() {
		return nil, errors.New("zero opening time")
	}

	if !endingTime.IsZero() && openingTime.After(endingTime) {
		return nil, errors.Errorf("opening time %s is after ending time %s", openingTime.Format(time.RFC3339), endingTime.Format(time.RFC3339))
	}

	if isClosed && endingTime.IsZero() {
		return nil, errors.New("zero ending time when period is closed")
	}

	return &Period{
		sobId:            sobId,
		periodId:         periodId,
		previousPeriodId: previousPeriodId,
		financialYear:    financialYear,
		number:           number,
		openingTime:      openingTime,
		endingTime:       endingTime,
		isClosed:         isClosed,
	}, nil
}

func (p Period) SobId() uuid.UUID {
	return p.sobId
}

func (p Period) PeriodId() uuid.UUID {
	return p.periodId
}

func (p Period) PreviousPeriodId() uuid.UUID {
	return p.previousPeriodId
}

func (p Period) FinancialYear() int {
	return p.financialYear
}

func (p Period) Number() int {
	return p.number
}

func (p Period) OpeningTime() time.Time {
	return p.openingTime
}

func (p Period) EndingTime() time.Time {
	return p.endingTime
}

func (p Period) IsClosed() bool {
	return p.isClosed
}
