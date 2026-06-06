package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// yearEndRetainedEarningsAccount is the raw account number for 本年利润.
// TODO: make this configurable per SoB in the future.
const yearEndRetainedEarningsAccount = "003103"

type PeriodPreCloseCheckHandler struct {
	readModel GeneralLedgerReadModel
}

func NewPeriodPreCloseCheckHandler(readModel GeneralLedgerReadModel) PeriodPreCloseCheckHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PeriodPreCloseCheckHandler{readModel: readModel}
}

func (h PeriodPreCloseCheckHandler) Handle(ctx context.Context, sobId, periodId uuid.UUID) (PreCloseCheck, error) {
	// Fetch period to know if this is a year-end period
	period, err := h.readModel.PeriodById(ctx, sobId, periodId)
	if err != nil {
		return PreCloseCheck{}, fmt.Errorf("failed to fetch period: %w", err)
	}

	// 1. Unposted journal check
	unpostedJournals, err := h.checkUnpostedJournals(ctx, sobId, periodId)
	if err != nil {
		return PreCloseCheck{}, fmt.Errorf("failed to check unposted journals: %w", err)
	}

	// If unposted journals check fails, skip remaining checks
	if unpostedJournals.Status != CheckStatusPassed {
		return PreCloseCheck{
			UnpostedJournals:         unpostedJournals,
			ProfitAndLossBalance:     PreCloseCheckPnLBalance{Status: CheckStatusUndetermined},
			CurrentYearProfitAccount: PreCloseCheckCurrentYearProfitAccount{Status: CheckStatusUndetermined},
			TrialBalance:             PreCloseCheckTrialBalance{Status: CheckStatusUndetermined},
		}, nil
	}

	// 2. P&L balance check
	pnlBalance, err := h.checkProfitAndLossBalance(ctx, sobId, periodId)
	if err != nil {
		return PreCloseCheck{}, fmt.Errorf("failed to check profit and loss balance: %w", err)
	}

	// If P&L balance check fails, skip remaining checks
	if pnlBalance.Status != CheckStatusPassed {
		return PreCloseCheck{
			UnpostedJournals:         unpostedJournals,
			ProfitAndLossBalance:     pnlBalance,
			CurrentYearProfitAccount: PreCloseCheckCurrentYearProfitAccount{Status: CheckStatusUndetermined},
			TrialBalance:             PreCloseCheckTrialBalance{Status: CheckStatusUndetermined},
		}, nil
	}

	// 3. Year-end check
	currentYearProfitAccount, err := h.checkCurrentYearProfitAccount(ctx, sobId, periodId, period.PeriodNumber)
	if err != nil {
		return PreCloseCheck{}, fmt.Errorf("failed to check year-end account: %w", err)
	}

	// If current year profit account check fails (and it's applicable), skip trial balance
	if currentYearProfitAccount.Status == CheckStatusFailed {
		return PreCloseCheck{
			UnpostedJournals:         unpostedJournals,
			ProfitAndLossBalance:     pnlBalance,
			CurrentYearProfitAccount: currentYearProfitAccount,
			TrialBalance:             PreCloseCheckTrialBalance{Status: CheckStatusUndetermined},
		}, nil
	}

	// 4. Trial balance check
	trialBalance, err := h.checkTrialBalance(ctx, sobId, periodId)
	if err != nil {
		return PreCloseCheck{}, fmt.Errorf("failed to check trial balance: %w", err)
	}

	return PreCloseCheck{
		UnpostedJournals:         unpostedJournals,
		ProfitAndLossBalance:     pnlBalance,
		TrialBalance:             trialBalance,
		CurrentYearProfitAccount: currentYearProfitAccount,
	}, nil
}

func (h PeriodPreCloseCheckHandler) checkUnpostedJournals(ctx context.Context, sobId, periodId uuid.UUID) (PreCloseCheckUnpostedJournals, error) {
	return checkUnpostedJournals(ctx, h.readModel, sobId, periodId)
}

func checkUnpostedJournals(ctx context.Context, readModel GeneralLedgerReadModel, sobId, periodId uuid.UUID) (PreCloseCheckUnpostedJournals, error) {
	page, size := 1, 3
	p, err := pageable.New(page, size)
	if err != nil {
		return PreCloseCheckUnpostedJournals{}, err
	}

	periodIdFilter, _ := filterable.NewFilter("periodId", filterable.OptEq, periodId.String())
	isPostedFilter, _ := filterable.NewFilter("isPosted", filterable.OptEq, false)
	combined := filterable.NewFilterable(
		filterable.TypeAND,
		filterable.NewFilterableAtom(periodIdFilter),
		filterable.NewFilterableAtom(isPostedFilter),
	)

	pageRequest := data.NewPageRequest(p, sortable.Unsorted(), combined)
	journalsPage, err := readModel.SearchJournals(ctx, sobId, pageRequest)
	if err != nil {
		return PreCloseCheckUnpostedJournals{}, err
	}

	count := journalsPage.NumberOfElements()
	journals := make([]PreCloseCheckJournal, 0, len(journalsPage.Content()))
	for _, j := range journalsPage.Content() {
		journals = append(journals, PreCloseCheckJournal{
			Id:              j.Id,
			DocumentNumber:  j.DocumentNumber,
			HeaderText:      j.HeaderText,
			Amount:          j.Amount,
			TransactionDate: j.TransactionDate,
			IsReviewed:      j.IsReviewed,
			IsAudited:       j.IsAudited,
		})
	}

	status := CheckStatusPassed
	if count != 0 {
		status = CheckStatusFailed
	}

	return PreCloseCheckUnpostedJournals{
		Status:   status,
		Count:    count,
		Journals: journals,
	}, nil
}

func (h PeriodPreCloseCheckHandler) checkProfitAndLossBalance(ctx context.Context, sobId, periodId uuid.UUID) (PreCloseCheckPnLBalance, error) {
	ledgers, err := h.readModel.ProfitAndLossLedgersHavingBalanceInPeriod(ctx, sobId, periodId)
	if err != nil {
		return PreCloseCheckPnLBalance{}, err
	}

	accounts := make([]PreCloseCheckPnLAccount, 0, len(ledgers))
	for _, l := range ledgers {
		accounts = append(accounts, PreCloseCheckPnLAccount{
			RawAccountNumber: l.Account.RawAccountNumber,
			AccountTitle:     l.Account.Title,
			EndingAmount:     l.EndingAmount,
		})
	}

	status := CheckStatusPassed
	if len(accounts) != 0 {
		status = CheckStatusFailed
	}

	return PreCloseCheckPnLBalance{
		Status:   status,
		Accounts: accounts,
	}, nil
}

func (h PeriodPreCloseCheckHandler) checkTrialBalance(ctx context.Context, sobId, periodId uuid.UUID) (PreCloseCheckTrialBalance, error) {
	return checkTrialBalance(ctx, h.readModel, sobId, periodId)
}

func checkTrialBalance(ctx context.Context, readModel GeneralLedgerReadModel, sobId, periodId uuid.UUID) (PreCloseCheckTrialBalance, error) {
	periodIdFilter, _ := filterable.NewFilter("periodId", filterable.OptEq, periodId.String())
	accountLevelFilter, _ := filterable.NewFilter("account.level", filterable.OptEq, 1)
	combined := filterable.NewFilterable(
		filterable.TypeAND,
		filterable.NewFilterableAtom(periodIdFilter),
		filterable.NewFilterableAtom(accountLevelFilter),
	)

	pageRequest := data.NewPageRequest(pageable.Unpaged(), sortable.Unsorted(), combined)
	ledgersPage, err := readModel.SearchLedgers(ctx, sobId, pageRequest)
	if err != nil {
		return PreCloseCheckTrialBalance{}, err
	}

	var totalOpening, totalPeriod, totalEnding decimal.Decimal
	for _, l := range ledgersPage.Content() {
		totalOpening = totalOpening.Add(l.OpeningAmount)
		totalPeriod = totalPeriod.Add(l.PeriodAmount)
		totalEnding = totalEnding.Add(l.EndingAmount)
	}

	status := CheckStatusPassed
	if !totalOpening.IsZero() || !totalPeriod.IsZero() || !totalEnding.IsZero() {
		status = CheckStatusFailed
	}

	return PreCloseCheckTrialBalance{
		Status:        status,
		OpeningAmount: totalOpening,
		PeriodAmount:  totalPeriod,
		EndingAmount:  totalEnding,
	}, nil
}

func (h PeriodPreCloseCheckHandler) checkCurrentYearProfitAccount(ctx context.Context, sobId, periodId uuid.UUID, periodNumber int) (PreCloseCheckCurrentYearProfitAccount, error) {
	if periodNumber != 12 {
		// Not applicable for non-year-end periods, so mark as passed
		return PreCloseCheckCurrentYearProfitAccount{Status: CheckStatusPassed}, nil
	}

	ledger, err := h.readModel.LedgerByRawAccountNumberInPeriod(ctx, sobId, yearEndRetainedEarningsAccount, periodId)
	if err != nil {
		return PreCloseCheckCurrentYearProfitAccount{}, err
	}
	if ledger == nil {
		// account not found or no ledger entry — treat as zero balance, which passes
		return PreCloseCheckCurrentYearProfitAccount{Status: CheckStatusPassed}, nil
	}

	status := CheckStatusPassed
	if !ledger.EndingAmount.IsZero() {
		status = CheckStatusFailed
	}

	return PreCloseCheckCurrentYearProfitAccount{
		Status:           status,
		RawAccountNumber: ledger.Account.RawAccountNumber,
		AccountTitle:     ledger.Account.Title,
		EndingAmount:     ledger.EndingAmount,
	}, nil
}
