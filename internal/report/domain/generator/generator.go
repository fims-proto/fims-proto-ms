package generator

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/general_ledger"
	"github/fims-proto/fims-proto-ms/internal/report/domain/general_ledger/balance_direction"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/formula_rule"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"
)

type Generator struct {
	r  *report.Report
	gl service.GeneralLedgerService

	periods      []*general_ledger.Period
	ledgersCache map[uuid.UUID][]*aggregatedLedger // in memory cache
}

func NewGenerator(report *report.Report, generalLedger service.GeneralLedgerService) *Generator {
	if report == nil {
		panic("nil report")
	}

	if generalLedger == nil {
		panic("nil general ledger repository")
	}

	return &Generator{
		r:  report,
		gl: generalLedger,

		periods:      nil,
		ledgersCache: make(map[uuid.UUID][]*aggregatedLedger),
	}
}

func (g *Generator) Report() *report.Report {
	return g.r
}

// Generate accept report template and instantiate it as report instance with period id, then generate the amounts of it.
func (g *Generator) Generate(
	ctx context.Context,
	newReportId,
	periodId uuid.UUID,
	title string,
	amountTypes []string,
) (*report.Report, error) {
	newReport, err := g.r.Instantiate(newReportId, periodId, title, amountTypes)
	if err != nil {
		return nil, err
	}

	// replace the report with new instance
	g.r = newReport
	if err = g.Regenerate(ctx); err != nil {
		return nil, err
	}

	return newReport, nil
}

// Regenerate calculates the report amounts as per the latest ledgers
func (g *Generator) Regenerate(ctx context.Context) error {
	if g.r.PeriodId() == uuid.Nil {
		return errors.NewSlugError("report-generate-emptyPeriod")
	}

	if err := g.preparePeriods(ctx); err != nil {
		return err
	}

	for _, section := range g.r.Sections() {
		if err := g.processSectionAmounts(ctx, section); err != nil {
			return err
		}
	}

	return nil
}

// summarise sections and items amounts to a section
func (g *Generator) processSectionAmounts(ctx context.Context, section *report.Section) error {
	amounts := make([]decimal.Decimal, len(g.r.AmountTypes()))

	// sum subsections amounts
	for _, subSection := range section.Sections() {
		if err := g.processSectionAmounts(ctx, subSection); err != nil {
			return err
		}
		for i := range g.r.AmountTypes() {
			amounts[i] = amounts[i].Add(subSection.Amounts()[i])
		}
	}

	// sum items amounts
	for _, item := range section.Items() {
		if err := g.processItemAmounts(ctx, item, amounts); err != nil {
			return err
		}
		for i := range g.r.AmountTypes() {
			amounts[i] = amounts[i].Add(item.Amounts()[i].Mul(decimal.NewFromInt(int64(item.SumFactor()))))
		}
	}

	section.SetAmounts(amounts)
	return nil
}

// collect amounts of the formulas of an item
func (g *Generator) processItemAmounts(ctx context.Context, item *report.Item, sums []decimal.Decimal) error {
	amounts := make([]decimal.Decimal, len(g.r.AmountTypes()))

	switch item.DataSource() {
	case data_source.Sum:
		// use items sum directly
		copy(amounts, sums)

	case data_source.Formulas:
		// calculate from the formulas
		for _, formula := range item.Formulas() {
			if err := g.processFormulaAmounts(ctx, formula); err != nil {
				return err
			}
			for i := range g.r.AmountTypes() {
				amounts[i] = amounts[i].Add(formula.Amounts()[i].Mul(decimal.NewFromInt(int64(formula.SumFactor()))))
			}
		}
	}

	item.SetAmounts(amounts)
	return nil
}

func (g *Generator) processFormulaAmounts(ctx context.Context, formula *report.Formula) error {
	// get aggregated ledger with balances for the account
	ledgerAmounts, err := g.aggregateLedgers(ctx, formula.AccountId())
	if err != nil {
		return err
	}

	// collect as per formula rule
	amounts := make([]decimal.Decimal, len(g.r.AmountTypes()))
	for i := range g.r.AmountTypes() {
		switch formula.Rule() {
		case formula_rule.Net:
			amounts[i] = ledgerAmounts[i].debit.Sub(ledgerAmounts[i].credit)
			if ledgerAmounts[i].direction == balance_direction.Credit {
				amounts[i] = amounts[i].Neg()
			}
		case formula_rule.Debit:
			amounts[i] = ledgerAmounts[i].debit
		case formula_rule.Credit:
			amounts[i] = ledgerAmounts[i].credit
		case formula_rule.Transaction:
			amounts[i] = decimal.Max(ledgerAmounts[i].debit, ledgerAmounts[i].credit)
		default:
			return fmt.Errorf("unsupported formula rule: %s", formula.Rule())
		}
	}

	formula.SetAmounts(amounts)
	return nil
}
