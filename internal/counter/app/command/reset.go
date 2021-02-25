package command

import (
	"context"
	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"

	"github.com/pkg/errors"
)

type CounterResetHandler struct {
	repo counter.Repository
}

func NewCounterResetHandler(repo counter.Repository) CounterResetHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CounterResetHandler{repo: repo}
}

func (h CounterResetHandler) Handle(ctx context.Context, cmd CounterResetCmd) error {
	return h.repo.UpdateCounter(
		ctx,
		cmd.CounterUUID,
		func(c *counter.Counter) (*counter.Counter, error) {
			if err := c.Reset(); err != nil {
				return nil, errors.Wrap(err, "counter reset handler unable to reset")
			}
			return c, nil
		},
	)
}
