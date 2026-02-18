package generator

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"testing"

	"github/fims-proto/fims-proto-ms/internal/report/domain/general_ledger"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/validator"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type ledgerSample struct {
	periodKey     string // fiscal year + number
	accountId     uuid.UUID
	direction     string
	openingDebit  int64
	openingCredit int64
	debit         int64
	credit        int64
	endingDebit   int64
	endingCredit  int64
}

func TestGenerator_Regenerate(t *testing.T) {
	// given
	r := prepareReport(t)
	noOpValidator := &validator.NoOpValidator{}
	g := NewGenerator(r, mockGeneralLedgerRepository{}, noOpValidator)

	// when
	err := g.Regenerate(context.Background())

	// then
	assert.NoError(t, err)
	assert.Equal(t, decimal.RequireFromString("3000"), g.r.Sections()[0].Sections()[0].Items()[0].Amounts()[0])
	assert.Equal(t, decimal.RequireFromString("1222"), g.r.Sections()[0].Sections()[0].Items()[0].Amounts()[1])
	// assert.Equal(t, decimal.RequireFromString("558"), g.r.Sections()[0].Sections()[0].Items()[0].Amounts()[2])
	// assert.Equal(t, decimal.RequireFromString("556"), g.r.Sections()[0].Sections()[0].Items()[0].Amounts()[3])

	assert.Equal(t, decimal.RequireFromString("3000"), g.r.Sections()[0].Sections()[0].Items()[1].Amounts()[0])
	assert.Equal(t, decimal.RequireFromString("1222"), g.r.Sections()[0].Sections()[0].Items()[1].Amounts()[1])
	// assert.Equal(t, decimal.RequireFromString("558"), g.r.Sections()[0].Sections()[0].Items()[1].Amounts()[2])
	// assert.Equal(t, decimal.RequireFromString("556"), g.r.Sections()[0].Sections()[0].Items()[1].Amounts()[3])

	assert.Equal(t, decimal.RequireFromString("-4000"), g.r.Sections()[0].Sections()[0].Items()[2].Amounts()[0])
	assert.Equal(t, decimal.RequireFromString("-2944"), g.r.Sections()[0].Sections()[0].Items()[2].Amounts()[1])
	// assert.Equal(t, decimal.RequireFromString("4057"), g.r.Sections()[0].Sections()[0].Items()[2].Amounts()[2])
	// assert.Equal(t, decimal.RequireFromString("-1"), g.r.Sections()[0].Sections()[0].Items()[2].Amounts()[3])

	assert.Equal(t, decimal.RequireFromString("-1000"), g.r.Sections()[0].Sections()[0].Amounts()[0])
	assert.Equal(t, decimal.RequireFromString("-1722"), g.r.Sections()[0].Sections()[0].Amounts()[1])
	// assert.Equal(t, decimal.RequireFromString("4615"), g.r.Sections()[0].Sections()[0].Amounts()[2])
	// assert.Equal(t, decimal.RequireFromString("555"), g.r.Sections()[0].Sections()[0].Amounts()[3])

	assert.Equal(t, decimal.RequireFromString("-1000"), g.r.Sections()[0].Amounts()[0])
	assert.Equal(t, decimal.RequireFromString("-1722"), g.r.Sections()[0].Amounts()[1])
	// assert.Equal(t, decimal.RequireFromString("4615"), g.r.Sections()[0].Amounts()[2])
	// assert.Equal(t, decimal.RequireFromString("555"), g.r.Sections()[0].Amounts()[3])

	assert.Equal(t, decimal.RequireFromString("0"), g.r.Sections()[1].Items()[0].Amounts()[0])
	assert.Equal(t, decimal.RequireFromString("0"), g.r.Sections()[1].Items()[0].Amounts()[1])
	// assert.Equal(t, decimal.RequireFromString("2000"), g.r.Sections()[1].Items()[0].Amounts()[2])
	// assert.Equal(t, decimal.RequireFromString("500"), g.r.Sections()[1].Items()[0].Amounts()[3])

	assert.Equal(t, decimal.RequireFromString("0"), g.r.Sections()[1].Items()[1].Amounts()[0])
	assert.Equal(t, decimal.RequireFromString("0"), g.r.Sections()[1].Items()[1].Amounts()[1])
	// assert.Equal(t, decimal.RequireFromString("2400"), g.r.Sections()[1].Items()[1].Amounts()[2])
	// assert.Equal(t, decimal.RequireFromString("600"), g.r.Sections()[1].Items()[1].Amounts()[3])

	assert.Equal(t, decimal.RequireFromString("0"), g.r.Sections()[1].Amounts()[0])
	assert.Equal(t, decimal.RequireFromString("0"), g.r.Sections()[1].Amounts()[1])
	// assert.Equal(t, decimal.RequireFromString("-400"), g.r.Sections()[1].Amounts()[2])
	// assert.Equal(t, decimal.RequireFromString("-100"), g.r.Sections()[1].Amounts()[3])
}

var (
	periodId      = uuid.New()
	accountId0000 = uuid.New()
	accountId0001 = uuid.New()
	accountId0020 = uuid.New()
	accountId0021 = uuid.New()
	accountId100  = uuid.New()
	accountId110  = uuid.New()
)

var ledgerSamples = []ledgerSample{
	{"2024 3", accountId0000, "debit", 1000, 0, 0, 111, 889, 0},
	{"2024 4", accountId0000, "debit", 889, 0, 222, 0, 1111, 0},
	{"2024 5", accountId0000, "debit", 1111, 0, 0, 333, 778, 0},
	{"2024 6", accountId0000, "debit", 778, 0, 444, 0, 1222, 0},

	{"2024 3", accountId0001, "debit", 2000, 0, 0, 112, 1888, 0},
	{"2024 4", accountId0001, "debit", 1888, 0, 224, 0, 2112, 0},
	{"2024 5", accountId0001, "debit", 2112, 0, 0, 3360, 0, 1248},
	{"2024 6", accountId0001, "debit", 0, 1248, 112, 0, 0, 1136},

	{"2024 3", accountId0020, "credit", 3000, 0, 0, 131, 2869, 0},
	{"2024 4", accountId0020, "credit", 2869, 0, 0, 1310, 1559, 0},
	{"2024 5", accountId0020, "credit", 1559, 0, 0, 2620, 0, 1061},
	{"2024 6", accountId0020, "credit", 0, 1061, 1, 0, 0, 1060},

	{"2024 3", accountId0021, "debit", 4000, 0, 1, 0, 4001, 0},
	{"2024 4", accountId0021, "debit", 4001, 0, 1, 0, 4002, 0},
	{"2024 5", accountId0021, "debit", 4002, 0, 1, 0, 4003, 0},
	{"2024 6", accountId0021, "debit", 4003, 0, 1, 0, 4004, 0},

	{"2024 3", accountId100, "credit", 0, 0, 500, 500, 0, 0},
	{"2024 4", accountId100, "credit", 0, 0, 500, 500, 0, 0},
	{"2024 5", accountId100, "credit", 0, 0, 500, 500, 0, 0},
	{"2024 6", accountId100, "credit", 0, 0, 500, 500, 0, 0},

	{"2024 3", accountId110, "credit", 0, 0, 600, 600, 0, 0},
	{"2024 4", accountId110, "credit", 0, 0, 600, 600, 0, 0},
	{"2024 5", accountId110, "credit", 0, 0, 600, 600, 0, 0},
	{"2024 6", accountId110, "credit", 0, 0, 600, 600, 0, 0},
}

func prepareReport(t *testing.T) *report.Report {
	formula0000 := prepareFormula(t, 1, accountId0000, 1, "net")
	formula0001 := prepareFormula(t, 2, accountId0001, 1, "debit")
	formula0020 := prepareFormula(t, 1, accountId0020, 1, "credit")
	formula0021 := prepareFormula(t, 2, accountId0021, -1, "net")
	formula100 := prepareFormula(t, 1, accountId100, 1, "debit")
	formula110 := prepareFormula(t, 1, accountId110, 1, "credit")

	item000, err := report.NewItem(uuid.New(), "item_000", 1, 1, "", 1, false, "formulas", []*report.Formula{formula0000, formula0001}, nil, false, false, false)
	assert.NoError(t, err)

	item001, err := report.NewItem(uuid.New(), "item_001", 1, 2, "", 0, false, "sum", nil, nil, false, false, false)
	assert.NoError(t, err)

	item002, err := report.NewItem(uuid.New(), "item_002", 1, 3, "", 1, false, "formulas", []*report.Formula{formula0020, formula0021}, nil, false, false, false)
	assert.NoError(t, err)

	item10, err := report.NewItem(uuid.New(), "item_10", 1, 1, "", 1, false, "formulas", []*report.Formula{formula100}, nil, false, false, false)
	assert.NoError(t, err)

	item11, err := report.NewItem(uuid.New(), "item_11", 1, 2, "", -1, false, "formulas", []*report.Formula{formula110}, nil, false, false, false)
	assert.NoError(t, err)

	section00, err := report.NewSection(
		uuid.New(),
		"section_00",
		1,
		"", // sectionType
		nil,
		nil,
		[]*report.Item{item000, item001, item002},
	)
	assert.NoError(t, err)

	section0, err := report.NewSection(
		uuid.New(),
		"section_A",
		1,
		"", // sectionType
		nil,
		[]*report.Section{section00},
		nil,
	)
	assert.NoError(t, err)

	section1, err := report.NewSection(
		uuid.New(),
		"section_B",
		2,
		"", // sectionType
		nil,
		nil,
		[]*report.Item{item10, item11},
	)
	assert.NoError(t, err)

	r, err := report.New(
		uuid.New(),
		uuid.New(),
		periodId,
		"test report",
		false,
		"balance_sheet",
		[]string{"year_opening_balance", "period_ending_balance"},
		[]*report.Section{section0, section1},
	)
	assert.NoError(t, err)

	return r
}

func prepareFormula(t *testing.T, sequence int, accountId uuid.UUID, sumFactor int, rule string) *report.Formula {
	formula, err := report.NewFormula(uuid.New(), sequence, accountId, sumFactor, rule, nil)
	assert.NoError(t, err)
	return formula
}

type mockGeneralLedgerRepository struct{}

func (m mockGeneralLedgerRepository) ReadPeriodById(
	_ context.Context,
	_ uuid.UUID,
	_ uuid.UUID,
) (*general_ledger.Period, error) {
	return general_ledger.NewPeriod(2024, 6), nil
}

func (m mockGeneralLedgerRepository) ReadFirstPeriodOfTheYear(
	_ context.Context,
	_ uuid.UUID,
	_ int,
) (*general_ledger.Period, error) {
	return general_ledger.NewPeriod(2024, 3), nil
}

func (m mockGeneralLedgerRepository) ReadLedgersByAccountAndPeriodsOrderByPeriod(
	_ context.Context,
	_ uuid.UUID,
	accountId uuid.UUID,
	periods []*general_ledger.Period,
) ([]*general_ledger.Ledger, error) {
	var result []*general_ledger.Ledger

	for _, period := range periods {
		for _, sample := range ledgerSamples {
			if sample.accountId == accountId && sample.periodKey == fmt.Sprint(
				period.FiscalYear(),
				period.PeriodNumber(),
			) {
				account, _ := general_ledger.NewAccount(sample.direction)
				result = append(
					result, general_ledger.NewLedger(
						account,
						period,
						decimal.NewFromInt(sample.openingDebit),
						decimal.NewFromInt(sample.openingCredit),
						decimal.NewFromInt(sample.debit),
						decimal.NewFromInt(sample.credit),
						decimal.NewFromInt(sample.endingDebit),
						decimal.NewFromInt(sample.endingCredit),
					),
				)
			}
		}
	}

	slices.SortFunc(
		result, func(e, e2 *general_ledger.Ledger) int {
			if n := cmp.Compare(e.Period().FiscalYear(), e2.Period().FiscalYear()); n != 0 {
				return n
			}
			return cmp.Compare(e.Period().PeriodNumber(), e2.Period().PeriodNumber())
		},
	)

	return result, nil
}

func (m mockGeneralLedgerRepository) ReadPeriodIdByFiscalYearAndNumber(
	ctx context.Context,
	sobId uuid.UUID,
	fiscalYear, number int,
) (uuid.UUID, error) {
	panic("implement me")
}

func (m mockGeneralLedgerRepository) ReadAccountIdsByNumbers(
	ctx context.Context,
	sobId uuid.UUID,
	accountNumbers []string,
) (map[string]uuid.UUID, error) {
	panic("implement me")
}
