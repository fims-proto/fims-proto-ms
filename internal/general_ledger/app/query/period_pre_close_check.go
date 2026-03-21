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
	unpostedJournals, err := h.checkUnpostedJournals(ctx, sobId, periodId)
	if err != nil {
		return PreCloseCheck{}, fmt.Errorf("failed to check unposted journals: %w", err)
	}

	pnlBalance, err := h.checkProfitAndLossBalance(ctx, sobId, periodId)
	if err != nil {
		return PreCloseCheck{}, fmt.Errorf("failed to check profit and loss balance: %w", err)
	}

	trialBalance, err := h.checkTrialBalance(ctx, sobId, periodId)
	if err != nil {
		return PreCloseCheck{}, fmt.Errorf("failed to check trial balance: %w", err)
	}

	return PreCloseCheck{
		UnpostedJournals:     unpostedJournals,
		ProfitAndLossBalance: pnlBalance,
		TrialBalance:         trialBalance,
	}, nil
}

func (h PeriodPreCloseCheckHandler) checkUnpostedJournals(ctx context.Context, sobId, periodId uuid.UUID) (PreCloseCheckUnpostedJournals, error) {
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
	journalsPage, err := h.readModel.SearchJournals(ctx, sobId, pageRequest)
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

	return PreCloseCheckUnpostedJournals{
		Passed:   count == 0,
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
			AccountNumber: l.Account.AccountNumber,
			AccountTitle:  l.Account.Title,
			EndingAmount:  l.EndingAmount,
		})
	}

	return PreCloseCheckPnLBalance{
		Passed:   len(accounts) == 0,
		Accounts: accounts,
	}, nil
}

func (h PeriodPreCloseCheckHandler) checkTrialBalance(ctx context.Context, sobId, periodId uuid.UUID) (PreCloseCheckTrialBalance, error) {
	periodIdFilter, _ := filterable.NewFilter("periodId", filterable.OptEq, periodId.String())
	accountLevelFilter, _ := filterable.NewFilter("account.level", filterable.OptEq, 1)
	combined := filterable.NewFilterable(
		filterable.TypeAND,
		filterable.NewFilterableAtom(periodIdFilter),
		filterable.NewFilterableAtom(accountLevelFilter),
	)

	pageRequest := data.NewPageRequest(pageable.Unpaged(), sortable.Unsorted(), combined)
	ledgersPage, err := h.readModel.SearchLedgers(ctx, sobId, pageRequest)
	if err != nil {
		return PreCloseCheckTrialBalance{}, err
	}

	var totalOpening, totalPeriod, totalEnding decimal.Decimal
	for _, l := range ledgersPage.Content() {
		totalOpening = totalOpening.Add(l.OpeningAmount)
		totalPeriod = totalPeriod.Add(l.PeriodAmount)
		totalEnding = totalEnding.Add(l.EndingAmount)
	}

	passed := totalOpening.IsZero() && totalPeriod.IsZero() && totalEnding.IsZero()

	return PreCloseCheckTrialBalance{
		Passed:        passed,
		OpeningAmount: totalOpening,
		PeriodAmount:  totalPeriod,
		EndingAmount:  totalEnding,
	}, nil
}
