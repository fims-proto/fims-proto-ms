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

func (h CounterDataloadHandler) Handle(ctx context.Context) error {
	counter, err := domain.NewCounter(uuid.New(), generalVoucher, "记", "号")
	if err != nil {
		return errors.Wrap(err, "failed to load counter")
	}
	return h.repo.CreateCounter(ctx, counter)
}
