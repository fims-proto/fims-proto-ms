package period

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Period struct {
	id               uuid.UUID
	sobId            uuid.UUID
	previousPeriodId uuid.UUID
	financialYear    int
	periodNumber     int
	openingTime      time.Time
	endingTime       time.Time
	isClosed         bool
}

func New(id, sobId, previousPeriodId uuid.UUID, financialYear, periodNumber int, openingTime, endingTime time.Time, isClosed bool) (*Period, error) {
	if sobId == uuid.Nil {
		return nil, errors.New("nil sob id")
	}

	if id == uuid.Nil {
		return nil, errors.New("nil period id")
	}

	if financialYear < 1970 || financialYear > 9999 {
		return nil, errors.Errorf("invalid financial year %d", financialYear)
	}

	if periodNumber < 1 {
		return nil, errors.Errorf("invalid period number %d", periodNumber)
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
		id:               id,
		sobId:            sobId,
		previousPeriodId: previousPeriodId,
		financialYear:    financialYear,
		periodNumber:     periodNumber,
		openingTime:      openingTime,
		endingTime:       endingTime,
		isClosed:         isClosed,
	}, nil
}

func (p Period) Id() uuid.UUID {
	return p.id
}

func (p Period) SobId() uuid.UUID {
	return p.sobId
}

func (p Period) PreviousPeriodId() uuid.UUID {
	return p.previousPeriodId
}

func (p Period) FinancialYear() int {
	return p.financialYear
}

func (p Period) PeriodNumber() int {
	return p.periodNumber
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
