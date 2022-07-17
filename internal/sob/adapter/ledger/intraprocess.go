package ledger

import (
	"context"

	"github.com/google/uuid"
	ledgerPort "github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"
)

type IntraProcessAdapter struct {
	ledgerInterface ledgerPort.LedgerInterface
}

func NewIntraProcessAdapter(ledgerInterface ledgerPort.LedgerInterface) IntraProcessAdapter {
	return IntraProcessAdapter{ledgerInterface: ledgerInterface}
}

func (i IntraProcessAdapter) InitializeFirstPeriod(ctx context.Context, sobId uuid.UUID, financialYear, number int) error {
	return i.ledgerInterface.InitializeFirstPeriod(ctx, sobId, financialYear, number)
}
