package intraprocess

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/app"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type VoucherInterface struct {
	app app.Application
}

func NewVoucherInterface(app app.Application) VoucherInterface {
	return VoucherInterface{app: app}
}

func (i VoucherInterface) CheckVoucherPosted(ctx context.Context, voucherUUID uuid.UUID) (bool, error) {
	v, err := i.app.Queries.ReadVouchers.HandleReadByUUID(ctx, voucherUUID)
	if err != nil {
		return false, errors.Wrapf(err, "failed to read voucher %s", voucherUUID)
	}
	return v.IsPosted, nil
}
