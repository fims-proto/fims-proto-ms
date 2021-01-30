package app

import "github/fims-proto/fims-proto-ms/internal/account/app/query"

type Queries struct {
	ValidateAccounts query.ValidateAccountsHandler
}

type Commands struct{}

type Application struct {
	Queries  Queries
	Commands Commands
}
