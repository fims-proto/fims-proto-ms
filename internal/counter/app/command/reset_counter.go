package command

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"
)

type CounterResetHandler struct {
	repo counter.Repository
}

func NewCounterResetHandler(repo counter.Repository) CounterResetHandler {
	if repo == nil {
		panic(("nil repo"))
	}
	return CounterResetHandler{
		repo: repo,
	}
}

func (h CounterResetHandler) Handle(ctx context.Context, cmd CounterResetCmd) error {
	return h.repo.ResetCounter(
		ctx,
		cmd.UUID,
	)
}
