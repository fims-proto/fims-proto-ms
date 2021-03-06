package command

import (
	"context"

	"github.com/google/uuid"
)

type AccountService interface {
	// account number with all its superior number
	readSuperiorNumbers(ctx context.Context, accountNumber string) ([]string, error)
}

type VoucherService interface {
	checkVoucherPosted(ctx context.Context, voucherUUID uuid.UUID) (bool, error)
}
