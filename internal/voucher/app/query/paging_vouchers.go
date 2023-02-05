package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/service"
)

type PagingVouchersHandler struct {
	readModel   VoucherReadModel
	userService service.UserService
}

func NewPagingVouchersHandler(readModel VoucherReadModel, userService service.UserService) PagingVouchersHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if userService == nil {
		panic("nil user service")
	}

	return PagingVouchersHandler{
		readModel:   readModel,
		userService: userService,
	}
}

func (h PagingVouchersHandler) Handle(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Voucher], error) {
	entriesPage, err := h.readModel.SearchVouchers(ctx, sobId, pageRequest)
	if err != nil {
		return nil, err
	}

	vouchers := entriesPage.Content()

	vouchers, err = enrichUserName(ctx, h.userService, vouchers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to enrich user in voucher")
	}

	return data.NewPage(vouchers, pageRequest, entriesPage.NumberOfElements())
}
