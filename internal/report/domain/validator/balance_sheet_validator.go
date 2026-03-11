package validator

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	sectiontype "github/fims-proto/fims-proto-ms/internal/report/domain/report/section_type"
)

// BalanceSheetValidator validates the balance sheet equation: Assets = Liabilities + Equity
type BalanceSheetValidator struct{}

func (v *BalanceSheetValidator) Validate(_ context.Context, r *report.Report) error {
	// Find required sections by type
	assetsSection := findSectionByType(r.Sections(), sectiontype.Assets)
	if assetsSection == nil {
		return errors.NewSlugError("report-validation-missingSectionType", "assets")
	}

	liabilitiesSection := findSectionByType(r.Sections(), sectiontype.Liabilities)
	if liabilitiesSection == nil {
		return errors.NewSlugError("report-validation-missingSectionType", "liabilities")
	}

	equitySection := findSectionByType(r.Sections(), sectiontype.Equity)
	if equitySection == nil {
		return errors.NewSlugError("report-validation-missingSectionType", "equity")
	}

	// Validate balance for each amount type column
	assets := assetsSection.Amounts()
	liabilities := liabilitiesSection.Amounts()
	equity := equitySection.Amounts()

	for i, amountType := range r.AmountTypes() {
		liabilitiesEquity := liabilities[i].Add(equity[i])

		if !assets[i].Equal(liabilitiesEquity) {
			difference := assets[i].Sub(liabilitiesEquity)
			return errors.NewSlugError(
				"report-balanceSheet-imbalance",
				amountType.String(),
				assets[i].String(),
				liabilitiesEquity.String(),
				difference.String(),
			)
		}
	}

	return nil
}
