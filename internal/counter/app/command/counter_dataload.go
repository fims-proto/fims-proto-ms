package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/counter/domain"

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

func (h CounterDataloadHandler) Handle(ctx context.Context, sob string) error {
	m, err := domain.NewMatcher("-", sob, generalVoucher)
	if err != nil {
		return errors.Wrapf(err, "create counter matcher failed: %s, %s", sob, generalVoucher)
	}

	counter, err := domain.NewCounter(uuid.New(), *m, "记", "号")
	if err != nil {
		return errors.Wrap(err, "create counter failed")
	}
	return h.repo.CreateCounter(ctx, counter)
}
