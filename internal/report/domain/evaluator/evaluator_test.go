package evaluator

import (
	"context"
	"testing"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestEvaluator_LedgerAccountsUseColumnContext(t *testing.T) {
	accountId := uuid.New()
	sobId := uuid.New()
	periodId := uuid.New()
	r := newTestReport(t, sobId, periodId,
		[]testColumn{
			{"期末余额", report.ColumnPeriodEndingBalance},
			{"年初余额", report.ColumnYearOpeningBalance},
		},
		[]*report.Row{
			newLedgerRow(t, "BS_CASH", "货币资金", 1, accountId, report.LedgerMeasureNet),
		},
	)

	gl := newMockGL(periodId)
	gl.ledgers[accountId] = []service.Ledger{
		{
			AccountId:        accountId,
			BalanceDirection: "debit",
			Period:           *gl.firstPeriod,
			OpeningAmount:    decimal.NewFromInt(100),
			EndingAmount:     decimal.NewFromInt(130),
		},
		{
			AccountId:        accountId,
			BalanceDirection: "debit",
			Period:           *gl.currentPeriod,
			OpeningAmount:    decimal.NewFromInt(130),
			EndingAmount:     decimal.NewFromInt(150),
		},
	}

	ev := New(r, gl)
	require.NoError(t, ev.Regenerate(context.Background()))

	amounts := r.Rows()[0].Amounts()
	require.True(t, decimal.NewFromInt(150).Equal(amounts[0]))
	require.True(t, decimal.NewFromInt(100).Equal(amounts[1]))
}

func TestEvaluator_ChildrenSumUsesDirectChildrenAndRowsExplicit(t *testing.T) {
	revenueId := uuid.New()
	expenseId := uuid.New()
	nestedId := uuid.New()
	sobId := uuid.New()
	periodId := uuid.New()
	r := newTestReport(t, sobId, periodId,
		[]testColumn{
			{"本年累计金额", report.ColumnYearToDateAmount},
			{"本月金额", report.ColumnPeriodAmount},
		},
		[]*report.Row{
			newChildrenSumRow(t, "NET", "净额",
				newLedgerRow(t, "REV", "收入", 1, revenueId, report.LedgerMeasureTransaction),
				newNoneRow(t, "GROUP", "Nested group", 1,
					newLedgerRow(t, "IGNORED", "Ignored grandchild", 1, nestedId, report.LedgerMeasureTransaction),
				),
				newLedgerRow(t, "EXP", "费用", -1, expenseId, report.LedgerMeasureTransaction),
			),
			newRowsExplicitRow(t, "COPY_NET", "净额复制", []report.RowReference{{RowCode: "NET", SumFactor: 1}}),
		},
	)

	gl := newMockGL(periodId)
	gl.ledgers[revenueId] = []service.Ledger{
		{AccountId: revenueId, BalanceDirection: "credit", Period: service.Period{FiscalYear: 2026, PeriodNumber: 1}, PeriodCredit: decimal.NewFromInt(30)},
		{AccountId: revenueId, BalanceDirection: "credit", Period: *gl.currentPeriod, PeriodCredit: decimal.NewFromInt(40)},
	}
	gl.ledgers[expenseId] = []service.Ledger{
		{AccountId: expenseId, BalanceDirection: "debit", Period: service.Period{FiscalYear: 2026, PeriodNumber: 1}, PeriodDebit: decimal.NewFromInt(10)},
		{AccountId: expenseId, BalanceDirection: "debit", Period: *gl.currentPeriod, PeriodDebit: decimal.NewFromInt(15)},
	}
	gl.ledgers[nestedId] = []service.Ledger{
		{AccountId: nestedId, BalanceDirection: "credit", Period: *gl.currentPeriod, PeriodCredit: decimal.NewFromInt(999)},
	}

	ev := New(r, gl)
	require.NoError(t, ev.Regenerate(context.Background()))

	net := r.Rows()[0].Amounts()
	require.True(t, decimal.NewFromInt(45).Equal(net[0]))
	require.True(t, decimal.NewFromInt(25).Equal(net[1]))
	require.Equal(t, net, r.Rows()[1].Amounts())
}

func TestEvaluator_CashFlowUsesAbsAndEndingCashFormula(t *testing.T) {
	opInId := uuid.New()
	opOutId := uuid.New()
	cashAccountId := uuid.New()
	sobId := uuid.New()
	periodId := uuid.New()
	r := newTestReport(t, sobId, periodId,
		[]testColumn{
			{"本年累计金额", report.ColumnYearToDateAmount},
			{"本月金额", report.ColumnPeriodAmount},
		},
		[]*report.Row{
			newChildrenSumRow(t, "NET_INCREASE", "现金净增加额",
				newCashFlowRow(t, "OP_IN", "流入", 1, opInId),
				newCashFlowRow(t, "OP_OUT", "流出", -1, opOutId),
			),
			newLedgerRow(t, "OPENING_CASH", "期初现金余额", 0, cashAccountId, report.LedgerMeasureOpeningBalance),
			newRowsExplicitRow(t, "ENDING_CASH", "期末现金余额", []report.RowReference{
				{RowCode: "OPENING_CASH", SumFactor: 1},
				{RowCode: "NET_INCREASE", SumFactor: 1},
			}),
		},
	)

	gl := newMockGL(periodId)
	gl.cashFlowAmounts[opInId] = decimal.NewFromInt(100)
	gl.cashFlowAmounts[opOutId] = decimal.NewFromInt(30)
	gl.ledgers[cashAccountId] = []service.Ledger{
		{AccountId: cashAccountId, BalanceDirection: "debit", Period: *gl.firstPeriod, OpeningAmount: decimal.NewFromInt(500)},
		{AccountId: cashAccountId, BalanceDirection: "debit", Period: *gl.currentPeriod, OpeningAmount: decimal.NewFromInt(540)},
	}

	ev := New(r, gl)
	require.NoError(t, ev.Regenerate(context.Background()))

	net := r.Rows()[0].Amounts()
	require.True(t, decimal.NewFromInt(70).Equal(net[0]))
	require.True(t, decimal.NewFromInt(70).Equal(net[1]))

	opening := r.Rows()[1].Amounts()
	require.True(t, decimal.NewFromInt(500).Equal(opening[0]))
	require.True(t, decimal.NewFromInt(540).Equal(opening[1]))

	ending := r.Rows()[2].Amounts()
	require.True(t, decimal.NewFromInt(570).Equal(ending[0]))
	require.True(t, decimal.NewFromInt(610).Equal(ending[1]))
}

func TestEvaluator_RowsExplicitMissingReference(t *testing.T) {
	sobId := uuid.New()
	periodId := uuid.New()
	r := newTestReport(t, sobId, periodId,
		[]testColumn{{"本月金额", report.ColumnPeriodAmount}},
		[]*report.Row{
			newRowsExplicitRow(t, "BROKEN", "Broken", []report.RowReference{{RowCode: "MISSING", SumFactor: 1}}),
		},
	)

	ev := New(r, newMockGL(periodId))
	require.Error(t, ev.Regenerate(context.Background()))
}

type testColumn struct {
	label     string
	valueType string
}

func newTestReport(t *testing.T, sobId uuid.UUID, periodId uuid.UUID, columns []testColumn, rows []*report.Row) *report.Report {
	t.Helper()
	var reportColumns []*report.Column
	for i, c := range columns {
		col, err := report.NewColumn(uuid.New(), c.label, c.valueType, i+1)
		require.NoError(t, err)
		reportColumns = append(reportColumns, col)
	}
	r, err := report.New(uuid.New(), sobId, periodId, "Report", false, report.ClassIncomeStatement, reportColumns, rows)
	require.NoError(t, err)
	return r
}

func newLedgerRow(t *testing.T, code string, text string, sumFactor int, accountId uuid.UUID, measure string) *report.Row {
	t.Helper()
	expr, err := report.NewExpression(uuid.New(), report.ExpressionLedgerAccounts, []report.LedgerAccountReference{
		{AccountId: accountId, SumFactor: 1, Measure: measure},
	}, nil, nil)
	require.NoError(t, err)
	row, err := report.NewRow(uuid.New(), code, text, 1, nil, true, sumFactor, true, true, false, expr, nil, nil)
	require.NoError(t, err)
	return row
}

func newCashFlowRow(t *testing.T, code string, text string, sumFactor int, itemId uuid.UUID) *report.Row {
	t.Helper()
	expr, err := report.NewExpression(uuid.New(), report.ExpressionCashFlowItems, nil, []report.CashFlowItemReference{
		{ItemId: itemId, SumFactor: 1},
	}, nil)
	require.NoError(t, err)
	row, err := report.NewRow(uuid.New(), code, text, 1, nil, true, sumFactor, true, true, false, expr, nil, nil)
	require.NoError(t, err)
	return row
}

func newChildrenSumRow(t *testing.T, code string, text string, rows ...*report.Row) *report.Row {
	t.Helper()
	expr, err := report.NewExpression(uuid.New(), report.ExpressionChildrenSum, nil, nil, nil)
	require.NoError(t, err)
	row, err := report.NewRow(uuid.New(), code, text, 1, nil, true, 1, true, true, true, expr, rows, nil)
	require.NoError(t, err)
	return row
}

func newRowsExplicitRow(t *testing.T, code string, text string, refs []report.RowReference) *report.Row {
	t.Helper()
	expr, err := report.NewExpression(uuid.New(), report.ExpressionRowsExplicit, nil, nil, refs)
	require.NoError(t, err)
	row, err := report.NewRow(uuid.New(), code, text, 1, nil, true, 0, true, true, false, expr, nil, nil)
	require.NoError(t, err)
	return row
}

func newNoneRow(t *testing.T, code string, text string, sumFactor int, rows ...*report.Row) *report.Row {
	t.Helper()
	expr, err := report.NewExpression(uuid.New(), report.ExpressionNone, nil, nil, nil)
	require.NoError(t, err)
	row, err := report.NewRow(uuid.New(), code, text, 1, nil, true, sumFactor, true, true, true, expr, rows, nil)
	require.NoError(t, err)
	return row
}

type mockGL struct {
	currentPeriod   *service.Period
	firstPeriod     *service.Period
	ledgers         map[uuid.UUID][]service.Ledger
	cashFlowAmounts map[uuid.UUID]decimal.Decimal
}

func newMockGL(periodId uuid.UUID) *mockGL {
	return &mockGL{
		currentPeriod:   &service.Period{Id: periodId, FiscalYear: 2026, PeriodNumber: 2},
		firstPeriod:     &service.Period{FiscalYear: 2026, PeriodNumber: 1},
		ledgers:         make(map[uuid.UUID][]service.Ledger),
		cashFlowAmounts: make(map[uuid.UUID]decimal.Decimal),
	}
}

func (m *mockGL) ReadPeriodIdByFiscalYearAndNumber(context.Context, uuid.UUID, int, int) (uuid.UUID, error) {
	return m.currentPeriod.Id, nil
}

func (m *mockGL) ReadPeriodById(context.Context, uuid.UUID, uuid.UUID) (*service.Period, error) {
	return m.currentPeriod, nil
}

func (m *mockGL) ReadFirstPeriodOfTheYear(context.Context, uuid.UUID, int) (*service.Period, error) {
	return m.firstPeriod, nil
}

func (m *mockGL) ReadAccountIdsByRawNumbers(context.Context, uuid.UUID, []string) (map[string]uuid.UUID, error) {
	return nil, nil
}

func (m *mockGL) ReadCashFlowItemIdsByCodes(context.Context, uuid.UUID, []string) (map[string]uuid.UUID, error) {
	return nil, nil
}

func (m *mockGL) ReadLedgersByAccountAndPeriodsOrderByPeriod(_ context.Context, _ uuid.UUID, accountId uuid.UUID, periods []*service.Period) ([]service.Ledger, error) {
	var result []service.Ledger
	for _, ledger := range m.ledgers[accountId] {
		for _, period := range periods {
			if ledger.Period.FiscalYear == period.FiscalYear && ledger.Period.PeriodNumber == period.PeriodNumber {
				result = append(result, ledger)
				break
			}
		}
	}
	return result, nil
}

func (m *mockGL) SumAbsJournalLineAmountsByCashFlowItemsAndPeriods(context.Context, uuid.UUID, []uuid.UUID, []*service.Period) (map[uuid.UUID]decimal.Decimal, error) {
	return m.cashFlowAmounts, nil
}
