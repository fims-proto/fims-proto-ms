package command

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"

	"github.com/google/uuid"
)

type CounterCreateHandler struct {
	repo counter.Repository
}

func NewCounterCreateHandler(repo counter.Repository) CounterCreateHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CounterCreateHandler{repo: repo}
}

func (h CounterCreateHandler) Handle(ctx context.Context, cmd CounterCreateCmd) error {
	counter, err := counter.NewCounter(uuid.New(), cmd.BusinessObject, cmd.Prefix, cmd.Sufix)
	if err != nil {
		return err
	}
	return h.repo.CreateCounter(ctx, counter)
}
