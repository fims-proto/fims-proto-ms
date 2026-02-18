package general_ledger

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	generalLedgerInterface intraprocess.GeneralLedgerInterface
}

func NewIntraProcessAdapter(generalLedgerInterface intraprocess.GeneralLedgerInterface) IntraProcessAdapter {
	return IntraProcessAdapter{generalLedgerInterface: generalLedgerInterface}
}

func (i IntraProcessAdapter) InitializeGeneralLedger(ctx context.Context, sobId uuid.UUID) error {
	return i.generalLedgerInterface.Initialize(ctx, command.InitializeCmd{SobId: sobId})
}
