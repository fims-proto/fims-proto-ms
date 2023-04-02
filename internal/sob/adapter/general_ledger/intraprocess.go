package general_ledger

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"

	generalLedgerPort "github/fims-proto/fims-proto-ms/internal/general_ledger/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	generalLedgerInterface generalLedgerPort.GeneralLedgerInterface
}

func NewIntraProcessAdapter(accountInterface generalLedgerPort.GeneralLedgerInterface) IntraProcessAdapter {
	return IntraProcessAdapter{generalLedgerInterface: accountInterface}
}

func (i IntraProcessAdapter) InitializeForSob(ctx context.Context, sobId uuid.UUID) error {
	return i.generalLedgerInterface.Initialize(ctx, command.InitializeCmd{SobId: sobId})
}
