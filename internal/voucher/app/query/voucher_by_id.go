package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/service"
)

type VoucherByIdHandler struct {
	readModel      VoucherReadModel
	accountService service.AccountService
	userService    service.UserService
}

func NewVoucherByIdHandler(
	readModel VoucherReadModel,
	accountService service.AccountService,
	userService service.UserService,
) VoucherByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if accountService == nil {
		panic("nil account service")
	}

	if userService == nil {
		panic("nil user service")
	}

	return VoucherByIdHandler{
		readModel:      readModel,
		accountService: accountService,
		userService:    userService,
	}
}

func (h VoucherByIdHandler) Handle(ctx context.Context, voucherId uuid.UUID) (Voucher, error) {
	v, err := h.readModel.VoucherById(ctx, voucherId)
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to read voucher")
	}

	singletonList, err := enrichLineItemAccountNumber(ctx, h.accountService, []Voucher{v})
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to enrich account number in voucher")
	}

	singletonList, err = enrichUserName(ctx, h.userService, singletonList)
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to enrich user in voucher")
	}

	singletonList, err = enrichPeriod(ctx, h.accountService, singletonList)
	if err != nil {
		return Voucher{}, errors.Wrap(err, "failed to enrich period in voucher")
	}

	return singletonList[0], nil
}
