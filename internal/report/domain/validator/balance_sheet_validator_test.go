package validator

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"
	sectiontype "github/fims-proto/fims-proto-ms/internal/report/domain/report/section_type"
)

func TestBalanceSheetValidator_Balanced(t *testing.T) {
	// Given: A perfectly balanced balance sheet
	r := prepareBalancedBalanceSheet(t)
	validator := &BalanceSheetValidator{}

	// When: Validating
	err := validator.Validate(context.Background(), r)

	// Then: Should pass
	assert.NoError(t, err)
}

func TestBalanceSheetValidator_Imbalance(t *testing.T) {
	// Given: An imbalanced balance sheet (modify equity to create imbalance)
	r := prepareBalancedBalanceSheet(t)

	equitySection := findSectionByType(r.Sections(), sectiontype.Equity)
	// Tamper with equity amounts to create imbalance
	equitySection.SetAmounts([]decimal.Decimal{
		decimal.NewFromInt(5000), // Changed from 3000
		decimal.NewFromInt(3000), // Changed from 2000
	})

	validator := &BalanceSheetValidator{}

	// When: Validating
	err := validator.Validate(context.Background(), r)

	// Then: Should fail with imbalance error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "report-balanceSheet-imbalance")
}

func TestBalanceSheetValidator_MissingAssetsSection(t *testing.T) {
	// Given: A balance sheet without assets section type
	r := prepareBalanceSheetWithoutSectionType(t, sectiontype.Assets)
	validator := &BalanceSheetValidator{}

	// When: Validating
	err := validator.Validate(context.Background(), r)

	// Then: Should fail with missing section type error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "report-validation-missingSectionType")
}

func TestBalanceSheetValidator_MissingLiabilitiesSection(t *testing.T) {
	// Given: A balance sheet without liabilities section type
	r := prepareBalanceSheetWithoutSectionType(t, sectiontype.Liabilities)
	validator := &BalanceSheetValidator{}

	// When: Validating
	err := validator.Validate(context.Background(), r)

	// Then: Should fail
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "report-validation-missingSectionType")
}

func TestBalanceSheetValidator_MissingEquitySection(t *testing.T) {
	// Given: A balance sheet without equity section type
	r := prepareBalanceSheetWithoutSectionType(t, sectiontype.Equity)
	validator := &BalanceSheetValidator{}

	// When: Validating
	err := validator.Validate(context.Background(), r)

	// Then: Should fail
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "report-validation-missingSectionType")
}

// Helper: Prepare a balanced balance sheet for testing
func prepareBalancedBalanceSheet(t *testing.T) *report.Report {
	t.Helper()

	// Assets = 10000, 8000
	assetsSection, err := report.NewSection(
		uuid.New(),
		"资产",
		1,
		"assets",
		[]decimal.Decimal{decimal.NewFromInt(10000), decimal.NewFromInt(8000)},
		nil,
		[]*report.Item{}, // Simplified: no items
	)
	assert.NoError(t, err)

	// Liabilities = 7000, 6000
	liabilitiesSection, err := report.NewSection(
		uuid.New(),
		"负债",
		1,
		"liabilities",
		[]decimal.Decimal{decimal.NewFromInt(7000), decimal.NewFromInt(6000)},
		nil,
		[]*report.Item{},
	)
	assert.NoError(t, err)

	// Equity = 3000, 2000 (Assets = Liabilities + Equity: 10000 = 7000 + 3000, 8000 = 6000 + 2000)
	equitySection, err := report.NewSection(
		uuid.New(),
		"所有者权益",
		2,
		"equity",
		[]decimal.Decimal{decimal.NewFromInt(3000), decimal.NewFromInt(2000)},
		nil,
		[]*report.Item{},
	)
	assert.NoError(t, err)

	containerSection, err := report.NewSection(
		uuid.New(),
		"负债及所有者权益",
		2,
		"",
		nil,
		[]*report.Section{liabilitiesSection, equitySection},
		nil,
	)
	assert.NoError(t, err)

	r, err := report.New(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		"资产负债表测试",
		false,
		"balance_sheet",
		[]string{amount_type.YearOpeningBalance.String(), amount_type.PeriodEndingBalance.String()},
		[]*report.Section{assetsSection, containerSection},
	)
	assert.NoError(t, err)

	return r
}

// Helper: Prepare a balance sheet without specific section type
func prepareBalanceSheetWithoutSectionType(t *testing.T, missingSectionType sectiontype.SectionType) *report.Report {
	t.Helper()

	assetsSectionType := "assets"
	liabilitiesSectionType := "liabilities"
	equitySectionType := "equity"

	// Remove the specified section type
	switch missingSectionType {
	case sectiontype.Assets:
		assetsSectionType = ""
	case sectiontype.Liabilities:
		liabilitiesSectionType = ""
	case sectiontype.Equity:
		equitySectionType = ""
	}

	assetsSection, _ := report.NewSection(
		uuid.New(),
		"资产",
		1,
		assetsSectionType,
		[]decimal.Decimal{decimal.NewFromInt(10000)},
		nil,
		[]*report.Item{},
	)

	liabilitiesSection, _ := report.NewSection(
		uuid.New(),
		"负债",
		1,
		liabilitiesSectionType,
		[]decimal.Decimal{decimal.NewFromInt(7000)},
		nil,
		[]*report.Item{},
	)

	equitySection, _ := report.NewSection(
		uuid.New(),
		"所有者权益",
		2,
		equitySectionType,
		[]decimal.Decimal{decimal.NewFromInt(3000)},
		nil,
		[]*report.Item{},
	)

	containerSection, _ := report.NewSection(
		uuid.New(),
		"负债及所有者权益",
		2,
		"",
		nil,
		[]*report.Section{liabilitiesSection, equitySection},
		nil,
	)

	r, _ := report.New(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		"资产负债表测试",
		false,
		"balance_sheet",
		[]string{amount_type.PeriodEndingBalance.String()},
		[]*report.Section{assetsSection, containerSection},
	)

	return r
}
