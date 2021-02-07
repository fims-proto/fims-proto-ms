package command

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"
)

type CounterNextHandler struct {
	repo counter.Repository
}

func NewCounterNextHandler(repo counter.Repository) CounterNextHandler{
	if repo == nil {
		panic("nil repo")
	}
	return CounterNextHandler{repo: repo}
}

func (h CounterNextHandler) Handle(ctx context.Context, cmd CounterNextCmd) (string, error){
	return h.repo.GetNextFromCounter(
		ctx, 
		cmd.UUID)
}