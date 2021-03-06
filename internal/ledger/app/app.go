package app

import "github/fims-proto/fims-proto-ms/internal/ledger/app/command"

type Queries struct{}

type Commands struct {
	UpdateLedgerBalanceHandler command.UpdateLedgerBalanceHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}
