package command

import (
	"context"

	"github.com/google/uuid"
)

type AccountService interface {
	// account number with all its superior number
	ReadSuperiorNumbers(ctx context.Context, sob, accountNumber string) ([]string, error)
}

type VoucherService interface {
	CheckVoucherPosted(ctx context.Context, voucherUUID uuid.UUID) (bool, error)
}
