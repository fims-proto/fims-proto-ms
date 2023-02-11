package service

import (
	"context"

	"github.com/google/uuid"
)

type AccountService interface {
	InitializeAccounts(ctx context.Context, sobId uuid.UUID) error
	InitializeFirstPeriod(ctx context.Context, sobId uuid.UUID, fiscalYear, number int) error
}
