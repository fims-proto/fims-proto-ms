package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain/ledger"
)

type UpdateLedgerBalanceCmd struct{}

type UpdateLedgerBalanceHandler struct {
	repo ledger.Repository
}

func NewUpdateLedgerBalanceHandler(repo ledger.Repository) UpdateLedgerBalanceHandler {
	if repo == nil {
		panic("nil repo")
	}
	return UpdateLedgerBalanceHandler{repo: repo}
}

func (h UpdateLedgerBalanceHandler) Handle(ctx context.Context, cmd UpdateLedgerBalanceCmd) error {
	return nil
}
