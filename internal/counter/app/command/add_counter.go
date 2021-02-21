package command

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"
)

type CounterAddHandler struct {
	repo counter.Repository
}

func NewCounterAddHandler(repo counter.Repository) CounterAddHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CounterAddHandler{repo: repo}
}

func (h CounterAddHandler) Handle(ctx context.Context, cmd CounterAddCmd) error {
	counter, err := counter.NewCounter(cmd.UUID, cmd.Length, cmd.Prefix, cmd.Sufix)
	if err != nil {
		return err
	}
	return h.repo.AddCounter(
		ctx,
		counter)
}
