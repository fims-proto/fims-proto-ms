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
	isCurrent    bool
}

// NewFuture creates valid period domain entity by given time point
// Typically used to create period when a specific time in future is provided when create a voucher
func NewFuture(id, sobId uuid.UUID, timePoint time.Time) (*Period, error) {
	fiscalYear, periodNumber := timePoint.Year(), int(timePoint.Month())
	openingTime := getOpeningTime(fiscalYear, periodNumber)
	if openingTime.Before(time.Now()) {
		return nil, errors.NewSlugError("period-timeInPast")
	}

	return NewByAllFields(
		id,
		sobId,
		fiscalYear,
		periodNumber,
		openingTime,
		getEndingTime(fiscalYear, periodNumber),
		false, // future period is always open
		false, // future period is not a current period
	)
}

// NewCurrent creates valid period domain entity by given fiscal year and number
// Typically used when initializing first period or closing and opening a new period
func NewCurrent(id, sobId uuid.UUID, fiscalYear, periodNumber int) (*Period, error) {
	return NewByAllFields(
		id,
		sobId,
		fiscalYear,
		periodNumber,
		getOpeningTime(fiscalYear, periodNumber),
		getEndingTime(fiscalYear, periodNumber),
		false, // current period is always open
		true,  // current period
	)
}

// NewByAllFields creates valid period domain entity by given all fields
// Typically used by other NewByXX methods or create from persistent entry
func NewByAllFields(
	id, sobId uuid.UUID,
	fiscalYear, periodNumber int,
	openingTime, endingTime time.Time,
	isClosed, isCurrent bool,
) (*Period, error) {
	if id == uuid.Nil {
		return nil, errors.NewSlugError("period-emptyId")
	}

	if sobId == uuid.Nil {
		return nil, errors.NewSlugError("emptySobId")
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
		isCurrent:    isCurrent,
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

func (p Period) IsCurrent() bool {
	return p.isCurrent
}
