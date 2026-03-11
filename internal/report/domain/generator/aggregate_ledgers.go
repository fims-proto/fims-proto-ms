package generator

import (
	"cmp"
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	"github/fims-proto/fims-proto-ms/internal/report/domain/general_ledger"
	"github/fims-proto/fims-proto-ms/internal/report/domain/general_ledger/balance_direction"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type aggregatedLedger struct {
	direction balance_direction.BalanceDirection
	amount    decimal.Decimal // signed: used for "Net" rule
	debit     decimal.Decimal // for "Debit" rule
	credit    decimal.Decimal // for "Credit" rule
}

func (g *Generator) preparePeriods(ctx context.Context) error {
	// determine the periods needed, as per the amount types
	periodsNeeded := make(map[string]*general_ledger.Period)
	periodKey := func(p *general_ledger.Period) string { return fmt.Sprint(p.FiscalYear(), p.PeriodNumber()) }

	// get report period
	currentPeriod, err := g.gl.ReadPeriodById(ctx, g.r.SobId(), g.r.PeriodId())
	if err != nil {
		return fmt.Errorf("failed to reading period: %w", err)
	}

	// get first period of the year or last year
	var firstPeriod *general_ledger.Period
	if slices.Contains(g.r.AmountTypes(), amount_type.LastYearAmount) {
		firstPeriod, err = g.gl.ReadFirstPeriodOfTheYear(ctx, g.r.SobId(), currentPeriod.FiscalYear()-1)
	}
	if firstPeriod == nil || slices.Contains(g.r.AmountTypes(), amount_type.YearToDateAmount) {
		firstPeriod, err = g.gl.ReadFirstPeriodOfTheYear(ctx, g.r.SobId(), currentPeriod.FiscalYear())
	}
	if firstPeriod == nil {
		return fmt.Errorf("first period not found of the year %d", currentPeriod.FiscalYear())
	}
	if err != nil {
		return fmt.Errorf("failed to read first period: %w", err)
	}

	for _, amountType := range g.r.AmountTypes() {
		switch amountType {
		case amount_type.YearOpeningBalance:
			periodsNeeded[periodKey(firstPeriod)] = firstPeriod
		case amount_type.PeriodEndingBalance:
			periodsNeeded[periodKey(currentPeriod)] = currentPeriod
		case amount_type.LastYearAmount, amount_type.YearToDateAmount:
			start := time.Date(firstPeriod.FiscalYear(), time.Month(firstPeriod.PeriodNumber()), 1, 0, 0, 0, 0, time.UTC)
			end := time.Date(currentPeriod.FiscalYear(), time.Month(currentPeriod.PeriodNumber()), 1, 0, 0, 0, 0, time.UTC)

			for current := start; !current.After(end); current = current.AddDate(0, 1, 0) {
				p := general_ledger.NewPeriod(current.Year(), int(current.Month()))
				periodsNeeded[periodKey(p)] = p
			}
		case amount_type.PeriodAmount:
			periodsNeeded[periodKey(currentPeriod)] = currentPeriod
		default:
			return fmt.Errorf("unsupported amount type: %s", amountType)
		}
	}

	periods := slices.Collect(maps.Values(periodsNeeded))
	sortByFiscalYearAndPeriod := func(a, b *general_ledger.Period) int {
		if n := cmp.Compare(a.FiscalYear(), b.FiscalYear()); n != 0 {
			return n
		}
		return cmp.Compare(a.PeriodNumber(), b.PeriodNumber())
	}
	slices.SortFunc(periods, sortByFiscalYearAndPeriod)

	g.periods = periods
	return nil
}

func (g *Generator) aggregateLedgers(ctx context.Context, accountId uuid.UUID) ([]*aggregatedLedger, error) {
	if cachedLedgers, ok := g.ledgersCache[accountId]; ok {
		// result is cached, use it directly
		return cachedLedgers, nil
	}

	// get ledgers from account id and periods
	ledgers, err := g.gl.ReadLedgersByAccountAndPeriodsOrderByPeriod(ctx, g.r.SobId(), accountId, g.periods)
	if err != nil {
		return nil, fmt.Errorf("failed to read ledgers: %w", err)
	}
	firstLedger := ledgers[0]
	currentLedger := ledgers[len(ledgers)-1]

	// aggregate ledgers as per amount types
	aggregated := make([]*aggregatedLedger, len(g.r.AmountTypes()))
	for i, amountType := range g.r.AmountTypes() {
		switch amountType {
		case amount_type.YearOpeningBalance:
			aggregated[i] = &aggregatedLedger{
				direction: firstLedger.Account().BalanceDirection(),
				amount:    firstLedger.OpeningAmount(),
				debit:     decimal.Zero, // Not used for opening balance
				credit:    decimal.Zero, // Not used for opening balance
			}
		case amount_type.PeriodEndingBalance:
			aggregated[i] = &aggregatedLedger{
				direction: currentLedger.Account().BalanceDirection(),
				amount:    currentLedger.EndingAmount(),
				debit:     decimal.Zero, // Not used for ending balance
				credit:    decimal.Zero, // Not used for ending balance
			}
		case amount_type.YearToDateAmount:
			// accumulated amount of the current year
			totalAmount := decimal.Zero
			totalDebit := decimal.Zero
			totalCredit := decimal.Zero

			for _, l := range ledgers {
				if l.Period().FiscalYear() == currentLedger.Period().FiscalYear() {
					totalAmount = totalAmount.Add(l.PeriodAmount())
					totalDebit = totalDebit.Add(l.PeriodDebit())
					totalCredit = totalCredit.Add(l.PeriodCredit())
				}
			}

			aggregated[i] = &aggregatedLedger{
				direction: firstLedger.Account().BalanceDirection(),
				amount:    totalAmount,
				debit:     totalDebit,
				credit:    totalCredit,
			}
		case amount_type.LastYearAmount:
			// accumulated amount of the previous year
			totalAmount := decimal.Zero
			totalDebit := decimal.Zero
			totalCredit := decimal.Zero

			for _, l := range ledgers {
				if l.Period().FiscalYear() == currentLedger.Period().FiscalYear()-1 {
					totalAmount = totalAmount.Add(l.PeriodAmount())
					totalDebit = totalDebit.Add(l.PeriodDebit())
					totalCredit = totalCredit.Add(l.PeriodCredit())
				}
			}

			aggregated[i] = &aggregatedLedger{
				direction: firstLedger.Account().BalanceDirection(),
				amount:    totalAmount,
				debit:     totalDebit,
				credit:    totalCredit,
			}
		case amount_type.PeriodAmount:
			aggregated[i] = &aggregatedLedger{
				direction: currentLedger.Account().BalanceDirection(),
				amount:    currentLedger.PeriodAmount(),
				debit:     currentLedger.PeriodDebit(),
				credit:    currentLedger.PeriodCredit(),
			}
		default:
			return nil, fmt.Errorf("unsupported amount type: %s", amountType)
		}
	}

	// cache the result
	g.ledgersCache[accountId] = aggregated
	return aggregated, nil
}
