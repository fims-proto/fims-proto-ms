package report

import (
	"slices"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/class"
)

type Report struct {
	id          uuid.UUID
	sobId       uuid.UUID
	periodId    uuid.UUID
	title       string
	template    bool
	class       class.Class
	amountTypes []amount_type.AmountType
	sections    []*Section
}

func New(
	id uuid.UUID,
	sobId uuid.UUID,
	periodId uuid.UUID,
	title string,
	template bool,
	reportClass string,
	amountTypes []string,
	sections []*Section,
) (*Report, error) {
	if id == uuid.Nil {
		return nil, errors.NewSlugError("report-emptyId")
	}

	if sobId == uuid.Nil {
		return nil, errors.NewSlugError("report-emptySobId")
	}

	if template && periodId != uuid.Nil {
		return nil, errors.NewSlugError("report-templateHasPeriod")
	}

	if !template && periodId == uuid.Nil {
		return nil, errors.NewSlugError("report-emptyPeriodId")
	}

	if title == "" {
		return nil, errors.NewSlugError("report-emptyTitle")
	}

	if len(sections) == 0 {
		return nil, errors.NewSlugError("report-emptySections")
	}

	newClass, err := class.FromString(reportClass)
	if err != nil {
		return nil, err
	}

	if len(amountTypes) == 0 {
		return nil, errors.NewSlugError("report-emptyAmountTypes")
	}

	var newAmountTypes []amount_type.AmountType
	for _, amountType := range amountTypes {
		newAmountType, err := amount_type.FromString(amountType)
		if err != nil {
			return nil, err
		}
		newAmountTypes = append(newAmountTypes, newAmountType)
	}

	if newClass == class.BalanceSheet && !containsOnly(
		newAmountTypes,
		amount_type.YearOpeningBalance,
		amount_type.PeriodEndingBalance,
	) {
		return nil, errors.NewSlugError("report-invalidAmountType")
	}

	if newClass == class.IncomeStatement && !containsOnly(
		newAmountTypes,
		amount_type.LastYearAmount,
		amount_type.YearToDateAmount,
		amount_type.PeriodAmount,
	) {
		return nil, errors.NewSlugError("report-invalidAmountType")
	}

	return &Report{
		id:          id,
		sobId:       sobId,
		periodId:    periodId,
		title:       title,
		template:    template,
		class:       newClass,
		amountTypes: newAmountTypes,
		sections:    sections,
	}, nil
}

// Instantiate deep copies report instance from template
func (r *Report) Instantiate(reportId uuid.UUID, periodId uuid.UUID) (*Report, error) {
	if !r.template {
		return nil, errors.NewSlugError("report-copyTemplate")
	}

	var newSections []*Section
	for _, section := range r.sections {
		newSections = append(newSections, section.copy())
	}

	var newAmountTypes []string
	for _, amountType := range r.amountTypes {
		newAmountTypes = append(newAmountTypes, amountType.String())
	}

	return New(
		reportId,
		r.sobId,
		periodId,
		r.title,
		false,
		r.class.String(),
		newAmountTypes,
		newSections,
	)
}

func (r *Report) Sections() []*Section {
	return r.sections
}

func (r *Report) PeriodId() uuid.UUID {
	return r.periodId
}

func (r *Report) AmountTypes() []amount_type.AmountType {
	return r.amountTypes
}

func (r *Report) Class() class.Class {
	return r.class
}

func (r *Report) Template() bool {
	return r.template
}

func (r *Report) Title() string {
	return r.title
}

func (r *Report) SobId() uuid.UUID {
	return r.sobId
}

func (r *Report) Id() uuid.UUID {
	return r.id
}

func containsOnly[E comparable](slice []E, targets ...E) bool {
	for _, element := range slice {
		if !slices.Contains(targets, element) {
			return false
		}
	}
	return true
}
