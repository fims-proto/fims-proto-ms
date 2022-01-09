package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/counter/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CounterCreateHandler struct {
	repo domain.Repository
}

func NewCounterCreateHandler(repo domain.Repository) CounterCreateHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CounterCreateHandler{repo: repo}
}

func (h CounterCreateHandler) Handle(ctx context.Context, cmd CounterCreateCmd) (err error) {
	log.Info(ctx, "handle creating counter")
	log.Debug(ctx, "handle creating counter, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle creating counter failed")
		}
	}()

	counter, err := domain.NewCounter(uuid.New(), 0, cmd.Prefix, cmd.Suffix, time.Time{}, ":", cmd.BusinessObjects...)
	if err != nil {
		return errors.Wrap(err, "create counter failed")
	}
	return h.repo.CreateCounter(ctx, counter)
}
