package voucher

import (
	"context"
	voucherport "github/fims-proto/fims-proto-ms/internal/voucher/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraprocessAdapter struct {
	voucherInterface voucherport.VoucherInterface
}

func NewIntraprocessAdapter(voucherInterface voucherport.VoucherInterface) IntraprocessAdapter {
	return IntraprocessAdapter{voucherInterface: voucherInterface}
}

func (s IntraprocessAdapter) CheckVoucherPosted(ctx context.Context, sob string, voucherUUID uuid.UUID) (bool, error) {
	return s.voucherInterface.CheckVoucherPosted(ctx, sob, voucherUUID)
}
