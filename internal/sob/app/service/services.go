package service

import (
	"context"

	"github.com/google/uuid"
)

type GeneralLedgerService interface {
	InitializeForSob(ctx context.Context, sobId uuid.UUID) error
}
