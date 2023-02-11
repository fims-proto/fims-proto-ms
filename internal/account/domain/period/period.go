package period

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

// Period consider a time point belongs to current period if ['openingTime', 'endingTime')
type Period struct {
	id           uuid.UUID
	sobId        uuid.UUID
	fiscalYear   int
	periodNumber int
	openingTime  time.Time
	endingTime   time.Time
	isClosed     bool
}

func New(id, sobId uuid.UUID, fiscalYear, periodNumber int, isClosed bool) (*Period, error) {
	if sobId == uuid.Nil {
		return nil, errors.NewSlugError("emptySobId")
	}

	if id == uuid.Nil {
		return nil, errors.NewSlugError("period-emptyId")
	}

	if fiscalYear < 1970 || fiscalYear > 9999 {
		return nil, errors.NewSlugError("period-invalidFiscalYear", fiscalYear)
	}

	if periodNumber < 1 || periodNumber > 12 {
		return nil, errors.NewSlugError("period-invalidPeriodNumber", periodNumber)
	}

	return NewFromPersistence(
		id,
		sobId,
		fiscalYear,
		periodNumber,
		getOpeningTime(fiscalYear, periodNumber),
		getEndingTime(fiscalYear, periodNumber),
		isClosed,
	)
}

func NewFromPersistence(id, sobId uuid.UUID, fiscalYear, periodNumber int, openingTime, endingTime time.Time, isClosed bool) (*Period, error) {
	if sobId == uuid.Nil {
		return nil, errors.NewSlugError("emptySobId")
	}

	if id == uuid.Nil {
		return nil, errors.NewSlugError("period-emptyId")
	}

	if fiscalYear < 1970 || fiscalYear > 9999 {
		return nil, errors.NewSlugError("period-invalidFiscalYear", fiscalYear)
	}

	if periodNumber < 1 || periodNumber > 12 {
		return nil, errors.NewSlugError("period-invalidPeriodNumber", periodNumber)
	}

	if openingTime.IsZero() {
		return nil, errors.NewSlugError("period-zeroOpeningTime")
	}

	if endingTime.IsZero() {
		return nil, errors.NewSlugError("period-zeroEndingTime")
	}

	expectedOpeningTime := getOpeningTime(fiscalYear, periodNumber)
	expectedEndingTime := getEndingTime(fiscalYear, periodNumber)

	if !openingTime.Equal(expectedOpeningTime) {
		return nil, errors.NewSlugError("period-invalidOpeningTime", openingTime, expectedOpeningTime)
	}

	if !endingTime.Equal(expectedEndingTime) {
		return nil, errors.NewSlugError("period-invalidEndingTime", endingTime, expectedEndingTime)
	}

	return &Period{
		id:           id,
		sobId:        sobId,
		fiscalYear:   fiscalYear,
		periodNumber: periodNumber,
		openingTime:  openingTime,
		endingTime:   endingTime,
		isClosed:     isClosed,
	}, nil
}

func (p Period) Id() uuid.UUID {
	return p.id
}

func (p Period) SobId() uuid.UUID {
	return p.sobId
}

func (p Period) FiscalYear() int {
	return p.fiscalYear
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
