package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
)

type initializePeriodCmd struct {
	SobId      uuid.UUID
	PeriodId   uuid.UUID
	FiscalYear int
	Number     int
}

func initializePeriod(ctx context.Context, cmd initializePeriodCmd, repo domain.Repository, numberingService service.NumberingService) error {
	return createPeriod(ctx, createPeriodCmd(cmd), repo, numberingService)
}
