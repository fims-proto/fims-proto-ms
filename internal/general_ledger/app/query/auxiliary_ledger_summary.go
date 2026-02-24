package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AuxiliaryLedgerSummary struct {
	AuxiliaryAccountId    uuid.UUID
	AuxiliaryAccountTitle string
	OpeningDebitBalance   decimal.Decimal
	OpeningCreditBalance  decimal.Decimal
	PeriodDebit           decimal.Decimal
	PeriodCredit          decimal.Decimal
	EndingDebitBalance    decimal.Decimal
	EndingCreditBalance   decimal.Decimal
}

type AuxiliaryLedgerSummaryHandler struct {
	readModel            GeneralLedgerReadModel
	periodRangeValidator periodRangeValidator
	categoryByKeyHandler AuxiliaryCategoryByKeyHandler
	accountByIdHandler   AccountByIdHandler
}

func NewAuxiliaryLedgerSummaryHandler(readModel GeneralLedgerReadModel) AuxiliaryLedgerSummaryHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return AuxiliaryLedgerSummaryHandler{
		readModel:            readModel,
		periodRangeValidator: newPeriodRangeValidator(readModel),
		categoryByKeyHandler: NewAuxiliaryCategoryByKeyHandler(readModel),
		accountByIdHandler:   NewAccountByIdHandler(readModel),
	}
}

func (h AuxiliaryLedgerSummaryHandler) Handle(
	ctx context.Context,
	sobId, accountId uuid.UUID,
	categoryKey, fromPeriod, toPeriod string,
	pageRequest data.PageRequest,
) (data.Page[AuxiliaryLedgerSummary], error) {
	// Get auxiliary category by key
	category, err := h.categoryByKeyHandler.Handle(ctx, sobId, categoryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get auxiliary category: %w", err)
	}

	// Get account to validate auxiliary category binding
	account, err := h.accountByIdHandler.Handle(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Validate that account is bound to the auxiliary category
	hasCategory := false
	for _, accCategory := range account.AuxiliaryCategories {
		if accCategory.Id == category.Id {
			hasCategory = true
			break
		}
	}
	if !hasCategory {
		return nil, errors.NewSlugError("account-auxiliary-not-bound", account.Title, category.Key)
	}

	// Validate period continuity
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber, err := h.periodRangeValidator.validate(ctx, sobId, fromPeriod, toPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid period range: %w", err)
	}

	// Query auxiliary ledgers for the period range with pagination
	ledgersPage, err := h.readModel.AuxiliariesByPeriodRange(
		ctx,
		sobId,
		accountId,
		category.Id,
		fromFiscalYear,
		fromPeriodNumber,
		toFiscalYear,
		toPeriodNumber,
		pageRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query auxiliary ledgers: %w", err)
	}

	// Aggregate balances per auxiliary account
	summaryMap := make(map[uuid.UUID]*AuxiliaryLedgerSummary)
	ledgers := ledgersPage.Content()

	for _, ledger := range ledgers {
		auxAccountId := ledger.AuxiliaryAccount.Id
		summary, exists := summaryMap[auxAccountId]

		if !exists {
			// First time seeing this auxiliary account - initialize with opening balance
			summary = &AuxiliaryLedgerSummary{
				AuxiliaryAccountId:    auxAccountId,
				AuxiliaryAccountTitle: ledger.AuxiliaryAccount.Title,
				OpeningDebitBalance:   ledger.OpeningDebitBalance,
				OpeningCreditBalance:  ledger.OpeningCreditBalance,
				PeriodDebit:           decimal.Zero,
				PeriodCredit:          decimal.Zero,
				EndingDebitBalance:    ledger.EndingDebitBalance,
				EndingCreditBalance:   ledger.EndingCreditBalance,
			}
			summaryMap[auxAccountId] = summary
		}

		// Accumulate period debit/credit across all periods
		summary.PeriodDebit = summary.PeriodDebit.Add(ledger.PeriodDebit)
		summary.PeriodCredit = summary.PeriodCredit.Add(ledger.PeriodCredit)
	}

	// Convert map to slice
	var summaries []AuxiliaryLedgerSummary
	for _, summary := range summaryMap {
		summaries = append(summaries, *summary)
	}

	// Return page with aggregated data
	return data.NewPage(
		summaries,
		pageRequest,
		len(summaries),
	)
}
