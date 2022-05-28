package command

import (
	"context"

	"github.com/google/uuid"
)

type AccountService interface {
	InitializeAccounts(ctx context.Context, sobId uuid.UUID) error
}

type CounterService interface {
	InitializeCounters(ctx context.Context, sobId uuid.UUID) error
}
