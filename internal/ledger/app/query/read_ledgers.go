package query

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LedgerReadModel interface {
	ReadLedgerById(ctx context.Context, id uuid.UUID) (Ledger, error)
	ReadAllLedgersByAccountingPeriod(ctx context.Context, periodId uuid.UUID) ([]Ledger, error)
	ReadAllAccountingPeriods(ctx context.Context, sobId uuid.UUID) ([]AccountingPeriod, error)
	ReadAccountingPeriodById(ctx context.Context, id uuid.UUID) (AccountingPeriod, error)
	ReadLedgerLogsByAccountIdsAndTimes(ctx context.Context, accountId []uuid.UUID, openingTime, endingTime time.Time) ([]LedgerLog, error)
}

type ReadLedgerHandler struct {
	readModel LedgerReadModel
}

func NewReadLedgerHandler(readModel LedgerReadModel) ReadLedgerHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return ReadLedgerHandler{readModel: readModel}
}

func (h ReadLedgerHandler) HandleReadAllAccountingPeriods(ctx context.Context, sobId uuid.UUID) ([]AccountingPeriod, error) {
	return h.readModel.ReadAllAccountingPeriods(ctx, sobId)
}

func (h ReadLedgerHandler) HandleReadAccountingPeriodById(ctx context.Context, id uuid.UUID) (AccountingPeriod, error) {
	return h.readModel.ReadAccountingPeriodById(ctx, id)
}

func (h ReadLedgerHandler) HandleReadAllLedgersByAccountingPeriod(ctx context.Context, periodId uuid.UUID) ([]Ledger, error) {
	return h.readModel.ReadAllLedgersByAccountingPeriod(ctx, periodId)
}
