package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github.com/google/uuid"
)

type AccountByIdHandler struct {
	readModel        GeneralLedgerReadModel
	dimensionService service.DimensionService
}

func NewAccountByIdHandler(readModel GeneralLedgerReadModel, dimensionService service.DimensionService) AccountByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if dimensionService == nil {
		panic("nil dimension service")
	}

	return AccountByIdHandler{
		readModel:        readModel,
		dimensionService: dimensionService,
	}
}

func (h AccountByIdHandler) Handle(ctx context.Context, accountId uuid.UUID) (Account, error) {
	idFilter, err := filterable.NewFilter("id", filterable.OptEq, accountId)
	if err != nil {
		panic("failed to create filter 'id'")
	}

	accounts, err := h.readModel.SearchAccounts(ctx, uuid.Nil, data.NewPageRequest(pageable.Unpaged(), sortable.Unsorted(), filterable.NewFilterableAtom(idFilter)))
	if err != nil {
		return Account{}, fmt.Errorf("failed to search accounts: %w", err)
	}
	if accounts.NumberOfElements() != 1 {
		return Account{}, commonErrors.ErrRecordNotFound()
	}

	account := accounts.Content()[0]

	account, err = enrichAccountDimensionCategories(ctx, h.dimensionService, account)
	if err != nil {
		return Account{}, fmt.Errorf("failed to enrich dimension categories: %w", err)
	}

	return account, nil
}
