package validator

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"
)

func TestIncomeStatementValidator_WithNetProfitItem(t *testing.T) {
	// Given: An income statement with a net profit item
	r := prepareIncomeStatementWithNetProfit(t)
	validator := &IncomeStatementValidator{}

	// When: Validating
	err := validator.Validate(context.Background(), r)

	// Then: Should pass
	assert.NoError(t, err)
}

func TestIncomeStatementValidator_MissingNetProfitItem(t *testing.T) {
	// Given: An income statement without a net profit item
	r := prepareIncomeStatementWithoutNetProfit(t)
	validator := &IncomeStatementValidator{}

	// When: Validating
	err := validator.Validate(context.Background(), r)

	// Then: Should fail
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "report-validation-missingItemType")
}

// Helper: Prepare an income statement with net profit item
func prepareIncomeStatementWithNetProfit(t *testing.T) *report.Report {
	t.Helper()

	netProfitItem, err := report.NewItem(uuid.New(), "净利润", 1, 1, "net_profit", 1, false, "sum", nil, []decimal.Decimal{decimal.NewFromInt(5000)}, false, false, false, false)
	assert.NoError(t, err)

	section, err := report.NewSection(
		uuid.New(),
		"利润",
		1,
		"",
		nil,
		nil,
		[]*report.Item{netProfitItem},
	)
	assert.NoError(t, err)

	r, err := report.New(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		"利润表测试",
		false,
		"income_statement",
		[]string{amount_type.PeriodAmount.String()},
		[]*report.Section{section},
	)
	assert.NoError(t, err)

	return r
}

// Helper: Prepare an income statement without net profit item
func prepareIncomeStatementWithoutNetProfit(t *testing.T) *report.Report {
	t.Helper()

	item, _ := report.NewItem(uuid.New(), "营业收入", 1, 1, "", 1, false, "sum", nil, []decimal.Decimal{decimal.NewFromInt(10000)}, false, false, false, false)

	section, _ := report.NewSection(
		uuid.New(),
		"收入",
		1,
		"",
		nil,
		nil,
		[]*report.Item{item},
	)

	r, _ := report.New(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		"利润表测试",
		false,
		"income_statement",
		[]string{amount_type.PeriodAmount.String()},
		[]*report.Section{section},
	)

	return r
}
