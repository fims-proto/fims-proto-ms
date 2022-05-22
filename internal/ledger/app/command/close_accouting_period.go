package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
)

// TODO

type CloseAccountingPeriodHandler struct {
	repo domain.Repository
}

func NewCloseAccountingPeriodHandler(repo domain.Repository) CloseAccountingPeriodHandler {
	if repo == nil {
		panic("nil ledger repo")
	}
	return CloseAccountingPeriodHandler{
		repo: repo,
	}
}

func (h CloseAccountingPeriodHandler) Handle(ctx context.Context) error {
	log.Info(ctx, "CloseAccountingPeriodHandler not implemented yet")
	panic("CloseAccountingPeriodHandler not implemented yet")
}
