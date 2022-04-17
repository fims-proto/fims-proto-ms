package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/log"

	counter "github/fims-proto/fims-proto-ms/internal/counter/domain"

	"github.com/pkg/errors"
)

type CounterNextHandler struct {
	repo counter.Repository
}

func NewCounterNextHandler(repo counter.Repository) CounterNextHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CounterNextHandler{repo: repo}
}

func (h CounterNextHandler) Handle(ctx context.Context, cmd CounterNextCmd) (identStr string, err error) {
	log.Info(ctx, "handle getting next counter")
	log.Debug(ctx, "handle getting next counter, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle getting next counter failed")
		}
	}()

	ident, err := h.repo.UpdateAndRead(
		ctx,
		cmd.CounterUUID,
		func(c *counter.Counter) (*counter.Counter, interface{}, error) {
			c.Next()
			return c, c.Identifier(), nil
		},
	)
	if err != nil {
		return "", errors.Wrap(err, "counter generate next identifier failed")
	}

	identStr, ok := ident.(string)
	if !ok {
		return "", errors.Errorf("expected identifier string, but got type %T", ident)
	}

	return identStr, nil
}
