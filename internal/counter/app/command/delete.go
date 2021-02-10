package command

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"
)

type CounterDeleteHandler struct {
	repo counter.Repository
}

func NewCounterDeleteHandler(repo counter.Repository) CounterDeleteHandler{
	if repo == nil {
		panic("nil repo")
	}
	return CounterDeleteHandler{repo: repo}
}

func (h CounterDeleteHandler) Handle(ctx context.Context, cmd CounterDeleteCmd) error{
	return h.repo.DeleteCounter(
		ctx, 
		cmd.UUID)
}