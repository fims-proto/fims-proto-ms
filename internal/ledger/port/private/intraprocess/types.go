package intraprocess

import (
	"github/fims-proto/fims-proto-ms/internal/ledger/app/command"
)

// use type from command directly, for now
type UpdateLedgerBalanceCmd struct {
	command.UpdateLedgerBalanceCmd
}
