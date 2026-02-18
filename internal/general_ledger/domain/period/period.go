package period

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

// Period represents an accounting period identified by fiscal year and period number.
// A period represents a calendar month when periodNumber 1-12 corresponds to January-December.
type Period struct {
	id           uuid.UUID
	sobId        uuid.UUID
	fiscalYear   int
	periodNumber int
	isClosed     bool
	isCurrent    bool
}

// New creates valid period domain entity by given fiscal year and number
// Typically used when initializing first period or closing and opening a new period
func New(id, sobId uuid.UUID, fiscalYear, periodNumber int, isCurrent bool) (*Period, error) {
	return NewByAllFields(
		id,
		sobId,
		fiscalYear,
		periodNumber,
		false, // never create a closed period
		isCurrent,
	)
}

// NewByAllFields creates valid period domain entity by given all fields
// Typically used by other NewByXX methods or create from persistent entry
func NewByAllFields(
	id uuid.UUID,
	sobId uuid.UUID,
	fiscalYear int,
	periodNumber int,
	isClosed bool,
	isCurrent bool,
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

	return &Period{
		id:           id,
		sobId:        sobId,
		fiscalYear:   fiscalYear,
		periodNumber: periodNumber,
		isClosed:     isClosed,
		isCurrent:    isCurrent,
	}, nil
}

func (p *Period) Id() uuid.UUID {
	return p.id
}

func (p *Period) SobId() uuid.UUID {
	return p.sobId
}

func (p *Period) FiscalYear() int {
	return p.fiscalYear
}

func (p *Period) PeriodNumber() int {
	return p.periodNumber
}

func (p *Period) IsClosed() bool {
	return p.isClosed
}

func (p *Period) IsCurrent() bool {
	return p.isCurrent
}

// NextNumber returns the fiscal year and period number of the next period.
// Handles year boundaries correctly (e.g., Period 12 of 2023 → Period 1 of 2024).
func (p *Period) NextNumber() (int, int) {
	firstDay := time.Date(p.fiscalYear, time.Month(p.periodNumber), 1, 0, 0, 0, 0, time.UTC)
	nextFirstDay := firstDay.AddDate(0, 1, 0)

	return nextFirstDay.Year(), int(nextFirstDay.Month())
}

// PreviousNumber returns the fiscal year and period number of the previous period.
// Handles year boundaries correctly (e.g., Period 1 of 2024 → Period 12 of 2023).
func (p *Period) PreviousNumber() (int, int) {
	firstDay := time.Date(p.fiscalYear, time.Month(p.periodNumber), 1, 0, 0, 0, 0, time.UTC)
	previousFirstDay := firstDay.AddDate(0, -1, 0)

	return previousFirstDay.Year(), int(previousFirstDay.Month())
}
