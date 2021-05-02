package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
)

type AccountService interface {
	ValidateExistence(ctx context.Context, accNumbers []string) error
}

type LedgerService interface {
	PostVoucher(ctx context.Context, voucher query.Voucher) error
}

type CounterService interface {
	GetNextIdentifier(ctx context.Context, businessObject string) (string, error)
}
