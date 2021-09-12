package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/counter/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	generalVoucher = string("GENERAL_VOUCHER")
)

type CounterDataloadHandler struct {
	repo domain.Repository
}

func NewCounterDataloadHandler(repo domain.Repository) CounterDataloadHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CounterDataloadHandler{repo: repo}
}

func (h CounterDataloadHandler) Handle(ctx context.Context, sob string) (err error) {
	log.Info(ctx, "handle counter dataload for sob %s", sob)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle counter dataload for sob %s failed", sob)
		}
	}()

	counter, err := domain.NewCounter(uuid.New(), 0, "记", "号", time.Time{}, "-", sob, generalVoucher)
	if err != nil {
		return errors.Wrap(err, "create counter failed")
	}
	return h.repo.CreateCounter(ctx, counter)
}
