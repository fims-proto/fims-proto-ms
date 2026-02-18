package validator

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/class"
	itemtype "github/fims-proto/fims-proto-ms/internal/report/domain/report/item_type"
	sectiontype "github/fims-proto/fims-proto-ms/internal/report/domain/report/section_type"
)

// ReportValidator validates report consistency based on its type
type ReportValidator interface {
	Validate(ctx context.Context, r *report.Report) error
}

// NoOpValidator does nothing, used for reports that don't need validation
type NoOpValidator struct{}

func (v *NoOpValidator) Validate(ctx context.Context, r *report.Report) error {
	return nil
}

// NewValidatorFactory creates the appropriate validator based on report class
func NewValidatorFactory(reportClass class.Class) ReportValidator {
	switch reportClass {
	case class.BalanceSheet:
		return &BalanceSheetValidator{}
	case class.IncomeStatement:
		return &IncomeStatementValidator{}
	default:
		return &NoOpValidator{}
	}
}

// findSectionByType recursively searches for a section with the specified type
func findSectionByType(sections []*report.Section, targetType sectiontype.SectionType) *report.Section {
	for _, section := range sections {
		if section.SectionType() == targetType {
			return section
		}

		// Recursively search in subsections
		if found := findSectionByType(section.Sections(), targetType); found != nil {
			return found
		}
	}
	return nil
}

// findItemByType recursively searches for an item with the specified type
func findItemByType(sections []*report.Section, targetType itemtype.ItemType) *report.Item {
	for _, section := range sections {
		// Search in items of this section
		for _, item := range section.Items() {
			if item.ItemType() == targetType {
				return item
			}
		}

		// Recursively search in subsections
		if found := findItemByType(section.Sections(), targetType); found != nil {
			return found
		}
	}
	return nil
}
