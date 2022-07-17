package self

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"
)

type IntraProcessAdapter struct {
	selfInterface intraprocess.LedgerInterface
}

func NewIntraProcessAdapter(selfInterface intraprocess.LedgerInterface) IntraProcessAdapter {
	return IntraProcessAdapter{selfInterface: selfInterface}
}

func (i IntraProcessAdapter) CreateLedgersForPeriod(ctx context.Context, periodId uuid.UUID) error {
	return i.selfInterface.InitializeLedgersForPeriod(ctx, periodId)
}
