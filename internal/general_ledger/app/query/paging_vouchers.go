package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
)

type PagingVouchersHandler struct {
	readModel   GeneralLedgerReadModel
	userService service.UserService
}

func NewPagingVouchersHandler(readModel GeneralLedgerReadModel, userService service.UserService) PagingVouchersHandler {
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
		return nil, fmt.Errorf("failed to enrich user in voucher: %w", err)
	}

	return data.NewPage(vouchers, pageRequest, entriesPage.NumberOfElements())
}
