package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
)

type LedgerService interface {
	LoadLedgers(ctx context.Context, accounts []domain.Account) error
}
