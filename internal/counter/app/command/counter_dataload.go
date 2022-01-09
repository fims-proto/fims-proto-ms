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

type CounterDataLoadHandler struct {
	repo domain.Repository
}

func NewCounterDataLoadHandler(repo domain.Repository) CounterDataLoadHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CounterDataLoadHandler{repo: repo}
}

func (h CounterDataLoadHandler) Handle(ctx context.Context, sobId uuid.UUID) (err error) {
	log.Info(ctx, "handle counter data load for sobId %s", sobId)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle counter data load for sobId %s failed", sobId)
		}
	}()

	counter, err := domain.NewCounter(uuid.New(), 0, "记", "号", time.Time{}, ":", sobId.String(), generalVoucher)
	if err != nil {
		return errors.Wrap(err, "create counter failed")
	}
	return h.repo.CreateCounter(ctx, counter)
}
