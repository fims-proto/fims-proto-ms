package service

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"

	"github.com/google/uuid"
)

type LedgerService interface {
	PostVoucher(ctx context.Context, voucher query.Voucher) error
}

type CounterService interface {
	GetNextIdentifier(ctx context.Context, businessObjects ...string) (string, error)
}

type AccountService interface {
	ValidateExistenceAndGetId(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]uuid.UUID, error)
}
