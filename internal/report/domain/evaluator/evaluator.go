package evaluator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Evaluator struct {
	r  *report.Report
	gl service.GeneralLedgerService

	currentPeriod *service.Period
	firstPeriod   *service.Period

	ledgerCache map[string][]service.Ledger
	rowAmounts  map[string][]decimal.Decimal
}

func New(report *report.Report, generalLedgerService service.GeneralLedgerService) *Evaluator {
	if report == nil {
		panic("nil report")
	}
	if generalLedgerService == nil {
		panic("nil general ledger service")
	}
	return &Evaluator{
		r:           report,
		gl:          generalLedgerService,
		ledgerCache: make(map[string][]service.Ledger),
		rowAmounts:  make(map[string][]decimal.Decimal),
	}
}

func (e *Evaluator) Report() *report.Report {
	return e.r
}

func (e *Evaluator) Regenerate(ctx context.Context) error {
	if e.r.PeriodId() == uuid.Nil {
		return fmt.Errorf("report period is required")
	}
	if err := e.preparePeriods(ctx); err != nil {
		return err
	}
	for _, row := range e.r.Rows() {
		if err := e.evaluateRowRecursive(ctx, row); err != nil {
			return err
		}
	}
	return nil
}

func (e *Evaluator) preparePeriods(ctx context.Context) error {
	current, err := e.gl.ReadPeriodById(ctx, e.r.SobId(), e.r.PeriodId())
	if err != nil {
		return fmt.Errorf("failed to read report period: %w", err)
	}
	first, err := e.gl.ReadFirstPeriodOfTheYear(ctx, e.r.SobId(), current.FiscalYear)
	if err != nil {
		return fmt.Errorf("failed to read first period of year: %w", err)
	}
	if first == nil {
		return fmt.Errorf("first period of year %d not found", current.FiscalYear)
	}
	e.currentPeriod = current
	e.firstPeriod = first
	return nil
}

func (e *Evaluator) evaluateRowRecursive(ctx context.Context, row *report.Row) error {
	childrenSum := zeroAmounts(len(e.r.Columns()))
	for _, child := range row.Rows() {
		if err := e.evaluateRowRecursive(ctx, child); err != nil {
			return err
		}
		childrenSum = addAmounts(childrenSum, child.Amounts(), child.SumFactor())
	}

	rowAmounts, err := e.evaluateRow(ctx, row, childrenSum)
	if err != nil {
		return err
	}
	row.SetAmounts(rowAmounts)
	e.rowAmounts[row.RowCode()] = rowAmounts
	return nil
}

func (e *Evaluator) evaluateRow(ctx context.Context, row *report.Row, childrenSum []decimal.Decimal) ([]decimal.Decimal, error) {
	switch row.Expression().Kind() {
	case report.ExpressionNone:
		return zeroAmounts(len(e.r.Columns())), nil
	case report.ExpressionChildrenSum:
		return append([]decimal.Decimal(nil), childrenSum...), nil
	case report.ExpressionRowsExplicit:
		return e.evaluateRowsExplicit(row)
	case report.ExpressionLedgerAccounts:
		return e.evaluateLedgerAccounts(ctx, row.Expression())
	case report.ExpressionCashFlowItems:
		return e.evaluateCashFlowItems(ctx, row.Expression())
	default:
		return nil, fmt.Errorf("unsupported expression kind %s", row.Expression().Kind())
	}
}

func (e *Evaluator) evaluateRowsExplicit(row *report.Row) ([]decimal.Decimal, error) {
	amounts := zeroAmounts(len(e.r.Columns()))
	for _, ref := range row.Expression().RowReferences() {
		refAmounts, ok := e.rowAmounts[ref.RowCode]
		if !ok {
			return nil, fmt.Errorf("row %s references missing or unevaluated row %s", row.RowCode(), ref.RowCode)
		}
		amounts = addAmounts(amounts, refAmounts, ref.SumFactor)
	}
	return amounts, nil
}

func (e *Evaluator) evaluateLedgerAccounts(ctx context.Context, expr *report.Expression) ([]decimal.Decimal, error) {
	amounts := zeroAmounts(len(e.r.Columns()))
	for _, accountRef := range expr.LedgerAccounts() {
		for i, column := range e.r.Columns() {
			value, err := e.evaluateLedgerAccount(ctx, accountRef, column.ValueType())
			if err != nil {
				return nil, err
			}
			amounts[i] = amounts[i].Add(value.Mul(decimal.NewFromInt(int64(accountRef.SumFactor))))
		}
	}
	return amounts, nil
}

func (e *Evaluator) evaluateLedgerAccount(ctx context.Context, accountRef report.LedgerAccountReference, columnType string) (decimal.Decimal, error) {
	periods, err := e.periodsForLedgerColumn(ctx, columnType, accountRef.Measure)
	if err != nil {
		return decimal.Zero, err
	}
	ledgers, err := e.readLedgers(ctx, accountRef.AccountId, periods)
	if err != nil {
		return decimal.Zero, err
	}
	if len(ledgers) == 0 {
		return decimal.Zero, fmt.Errorf("ledger not found for account %s", accountRef.AccountId)
	}

	switch columnType {
	case report.ColumnYearOpeningBalance:
		return e.applyLedgerMeasure(ledgers[0], accountRef.Measure, ledgerAmountOpening)
	case report.ColumnPeriodEndingBalance:
		return e.applyLedgerMeasure(ledgers[len(ledgers)-1], accountRef.Measure, ledgerAmountEnding)
	case report.ColumnPeriodAmount:
		if accountRef.Measure == report.LedgerMeasureOpeningBalance {
			return e.applyLedgerMeasure(ledgers[len(ledgers)-1], accountRef.Measure, ledgerAmountOpening)
		}
		return e.aggregateLedgers(ledgers, accountRef.Measure)
	case report.ColumnYearToDateAmount:
		if accountRef.Measure == report.LedgerMeasureOpeningBalance {
			return e.applyLedgerMeasure(ledgers[0], accountRef.Measure, ledgerAmountOpening)
		}
		return e.aggregateLedgers(ledgers, accountRef.Measure)
	case report.ColumnLastYearAmount:
		return e.aggregateLedgers(ledgers, accountRef.Measure)
	default:
		return decimal.Zero, fmt.Errorf("unsupported column type %s", columnType)
	}
}

type ledgerAmountPoint string

const (
	ledgerAmountOpening ledgerAmountPoint = "opening"
	ledgerAmountEnding  ledgerAmountPoint = "ending"
	ledgerAmountPeriod  ledgerAmountPoint = "period"
)

func (e *Evaluator) aggregateLedgers(ledgers []service.Ledger, measure string) (decimal.Decimal, error) {
	totalAmount := decimal.Zero
	totalDebit := decimal.Zero
	totalCredit := decimal.Zero
	for _, ledger := range ledgers {
		totalAmount = totalAmount.Add(ledger.PeriodAmount)
		totalDebit = totalDebit.Add(ledger.PeriodDebit)
		totalCredit = totalCredit.Add(ledger.PeriodCredit)
	}
	aggregated := service.Ledger{
		BalanceDirection: ledgers[0].BalanceDirection,
		PeriodAmount:     totalAmount,
		PeriodDebit:      totalDebit,
		PeriodCredit:     totalCredit,
	}
	return e.applyLedgerMeasure(aggregated, measure, ledgerAmountPeriod)
}

func (e *Evaluator) applyLedgerMeasure(ledger service.Ledger, measure string, point ledgerAmountPoint) (decimal.Decimal, error) {
	switch measure {
	case report.LedgerMeasureNet:
		var amount decimal.Decimal
		switch point {
		case ledgerAmountOpening:
			amount = ledger.OpeningAmount
		case ledgerAmountEnding:
			amount = ledger.EndingAmount
		default:
			amount = ledger.PeriodAmount
		}
		if ledger.BalanceDirection == "credit" {
			amount = amount.Neg()
		}
		return amount, nil
	case report.LedgerMeasureOpeningBalance:
		amount := ledger.OpeningAmount
		if ledger.BalanceDirection == "credit" {
			amount = amount.Neg()
		}
		return amount, nil
	case report.LedgerMeasureDebit:
		return ledger.PeriodDebit, nil
	case report.LedgerMeasureCredit:
		return ledger.PeriodCredit, nil
	case report.LedgerMeasureTransaction:
		return decimal.Max(ledger.PeriodDebit, ledger.PeriodCredit), nil
	default:
		return decimal.Zero, fmt.Errorf("unsupported ledger measure %s", measure)
	}
}

func (e *Evaluator) evaluateCashFlowItems(ctx context.Context, expr *report.Expression) ([]decimal.Decimal, error) {
	amounts := zeroAmounts(len(e.r.Columns()))
	for i, column := range e.r.Columns() {
		periods, err := e.periodsForCashFlowColumn(column.ValueType())
		if err != nil {
			return nil, err
		}
		var itemIds []uuid.UUID
		for _, ref := range expr.CashFlowItems() {
			itemIds = append(itemIds, ref.ItemId)
		}
		sums, err := e.gl.SumAbsJournalLineAmountsByCashFlowItemsAndPeriods(ctx, e.r.SobId(), itemIds, periods)
		if err != nil {
			return nil, fmt.Errorf("failed to aggregate cash flow item amounts: %w", err)
		}
		for _, ref := range expr.CashFlowItems() {
			amounts[i] = amounts[i].Add(sums[ref.ItemId].Mul(decimal.NewFromInt(int64(ref.SumFactor))))
		}
	}
	return amounts, nil
}

func (e *Evaluator) periodsForLedgerColumn(ctx context.Context, columnType string, measure string) ([]*service.Period, error) {
	switch columnType {
	case report.ColumnPeriodEndingBalance:
		return []*service.Period{e.currentPeriod}, nil
	case report.ColumnYearOpeningBalance:
		return []*service.Period{e.firstPeriod}, nil
	case report.ColumnPeriodAmount:
		return []*service.Period{e.currentPeriod}, nil
	case report.ColumnYearToDateAmount:
		if measure == report.LedgerMeasureOpeningBalance {
			return []*service.Period{e.firstPeriod}, nil
		}
		return e.periodRange(e.firstPeriod, e.currentPeriod), nil
	case report.ColumnLastYearAmount:
		first, err := e.gl.ReadFirstPeriodOfTheYear(ctx, e.r.SobId(), e.currentPeriod.FiscalYear-1)
		if err != nil {
			return nil, err
		}
		if first == nil {
			return nil, fmt.Errorf("first period of previous year not found")
		}
		last := &service.Period{FiscalYear: e.currentPeriod.FiscalYear - 1, PeriodNumber: e.currentPeriod.PeriodNumber}
		return e.periodRange(first, last), nil
	default:
		return nil, fmt.Errorf("unsupported column type %s", columnType)
	}
}

func (e *Evaluator) periodsForCashFlowColumn(columnType string) ([]*service.Period, error) {
	switch columnType {
	case report.ColumnPeriodAmount:
		return []*service.Period{e.currentPeriod}, nil
	case report.ColumnYearToDateAmount:
		return e.periodRange(e.firstPeriod, e.currentPeriod), nil
	default:
		return nil, fmt.Errorf("unsupported cash flow column type %s", columnType)
	}
}

func (e *Evaluator) periodRange(start *service.Period, end *service.Period) []*service.Period {
	var periods []*service.Period
	startDate := time.Date(start.FiscalYear, time.Month(start.PeriodNumber), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(end.FiscalYear, time.Month(end.PeriodNumber), 1, 0, 0, 0, 0, time.UTC)
	for current := startDate; !current.After(endDate); current = current.AddDate(0, 1, 0) {
		periods = append(periods, &service.Period{
			FiscalYear:   current.Year(),
			PeriodNumber: int(current.Month()),
		})
	}
	return periods
}

func (e *Evaluator) readLedgers(ctx context.Context, accountId uuid.UUID, periods []*service.Period) ([]service.Ledger, error) {
	cacheKey := accountId.String() + ":" + periodCacheKey(periods)
	if cached, ok := e.ledgerCache[cacheKey]; ok {
		return cached, nil
	}
	ledgers, err := e.gl.ReadLedgersByAccountAndPeriodsOrderByPeriod(ctx, e.r.SobId(), accountId, periods)
	if err != nil {
		return nil, fmt.Errorf("failed to read ledgers: %w", err)
	}
	e.ledgerCache[cacheKey] = ledgers
	return ledgers, nil
}

func periodCacheKey(periods []*service.Period) string {
	var parts []string
	for _, period := range periods {
		parts = append(parts, fmt.Sprintf("%d-%02d", period.FiscalYear, period.PeriodNumber))
	}
	return strings.Join(parts, ",")
}

func zeroAmounts(n int) []decimal.Decimal {
	amounts := make([]decimal.Decimal, n)
	for i := range amounts {
		amounts[i] = decimal.Zero
	}
	return amounts
}

func addAmounts(base []decimal.Decimal, addend []decimal.Decimal, factor int) []decimal.Decimal {
	if len(addend) == 0 {
		return base
	}
	f := decimal.NewFromInt(int64(factor))
	for i := range base {
		base[i] = base[i].Add(addend[i].Mul(f))
	}
	return base
}
