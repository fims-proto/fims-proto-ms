package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
)

type VoucherByIdHandler struct {
	readModel   GeneralLedgerReadModel
	userService service.UserService
}

func NewVoucherByIdHandler(
	readModel GeneralLedgerReadModel,
	userService service.UserService,
) VoucherByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if userService == nil {
		panic("nil user service")
	}

	return VoucherByIdHandler{
		readModel:   readModel,
		userService: userService,
	}
}

func (h VoucherByIdHandler) Handle(ctx context.Context, voucherId uuid.UUID) (Voucher, error) {
	v, err := h.readModel.VoucherById(ctx, voucherId)
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to read voucher")
	}

	singletonList, err := enrichUserName(ctx, h.userService, []Voucher{v})
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to enrich user in voucher")
	}

	return singletonList[0], nil
}
