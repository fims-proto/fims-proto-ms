package app

import "github/fims-proto/fims-proto-ms/internal/ledger/app/command"

type Queries struct{}

type Commands struct {
	UpdateLedgerBalance command.UpdateLedgerBalanceHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}
