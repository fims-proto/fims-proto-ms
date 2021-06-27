package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/counter/domain"
	"strings"

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

func (h CounterCreateHandler) Handle(ctx context.Context, cmd CounterCreateCmd) error {
	m, err := domain.NewMatcher("-", cmd.BusinessObjects...)
	if err != nil {
		return errors.Wrapf(err, "create counter matcher failed: %s", strings.Join(cmd.BusinessObjects, ","))
	}

	counter, err := domain.NewCounter(uuid.New(), *m, cmd.Prefix, cmd.Sufix)
	if err != nil {
		return errors.Wrap(err, "create counter failed")
	}
	return h.repo.CreateCounter(ctx, counter)
}
