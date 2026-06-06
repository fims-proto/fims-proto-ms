package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"

	"github.com/google/uuid"
)

type ClosingJournalIds struct {
	MonthlyClosingJournalId *uuid.UUID
	YearEndClosingJournalId *uuid.UUID
}

type ClosingJournalIdsByPeriodHandler struct {
	readModel GeneralLedgerReadModel
}

func NewClosingJournalIdsByPeriodHandler(readModel GeneralLedgerReadModel) ClosingJournalIdsByPeriodHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return ClosingJournalIdsByPeriodHandler{
		readModel: readModel,
	}
}

func (h ClosingJournalIdsByPeriodHandler) Handle(ctx context.Context, sobId uuid.UUID, fiscalYear, periodNumber int) (ClosingJournalIds, error) {
	monthly, err := h.readModel.ClosingJournalIdBySobAndPeriod(ctx, sobId, fiscalYear, periodNumber, string(journal.TypeClosing))
	if err != nil {
		return ClosingJournalIds{}, fmt.Errorf("failed to read monthly closing journal: %w", err)
	}

	yearEnd, err := h.readModel.ClosingJournalIdBySobAndPeriod(ctx, sobId, fiscalYear, periodNumber, string(journal.TypeYearlyClosing))
	if err != nil {
		return ClosingJournalIds{}, fmt.Errorf("failed to read year-end closing journal: %w", err)
	}

	return ClosingJournalIds{
		MonthlyClosingJournalId: monthly,
		YearEndClosingJournalId: yearEnd,
	}, nil
}
