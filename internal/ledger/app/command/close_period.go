package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
)

// TODO

type ClosePeriodHandler struct {
	repo domain.Repository
}

func NewClosePeriodHandler(repo domain.Repository) ClosePeriodHandler {
	if repo == nil {
		panic("nil ledger repo")
	}
	return ClosePeriodHandler{
		repo: repo,
	}
}

func (h ClosePeriodHandler) Handle(ctx context.Context) error {
	log.Info(ctx, "ClosePeriodHandler not implemented yet")
	panic("ClosePeriodHandler not implemented yet")
}
