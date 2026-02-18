package validator

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	itemtype "github/fims-proto/fims-proto-ms/internal/report/domain/report/item_type"
)

// IncomeStatementValidator validates the income statement structure and key items
type IncomeStatementValidator struct{}

func (v *IncomeStatementValidator) Validate(ctx context.Context, r *report.Report) error {
	// Find required net profit item
	netProfitItem := findItemByType(r.Sections(), itemtype.NetProfit)
	if netProfitItem == nil {
		return errors.NewSlugError("report-validation-missingItemType", "net_profit")
	}

	// For now, we only validate that the net profit item exists and has amounts
	// More complex validation (revenue - expenses = net profit) can be added later
	// when we have clear business rules for identifying revenue and expense items

	amounts := netProfitItem.Amounts()
	if len(amounts) == 0 {
		return errors.NewSlugError("report-incomeStatement-profitMismatch",
			"net_profit", "has values", "empty")
	}

	return nil
}
