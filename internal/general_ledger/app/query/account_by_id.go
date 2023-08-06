package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

type AccountByIdHandler struct {
	readModel GeneralLedgerReadModel
}

func NewAccountByIdHandler(readModel GeneralLedgerReadModel) AccountByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return AccountByIdHandler{
		readModel: readModel,
	}
}

func (h AccountByIdHandler) Handle(ctx context.Context, accountId uuid.UUID) (Account, error) {
	idFilter, err := filterable.NewFilter("id", "eq", accountId)
	if err != nil {
		return Account{}, fmt.Errorf("failed to build filter: %w", err)
	}

	accounts, err := h.readModel.SearchAccounts(ctx, uuid.Nil, data.NewPageRequest(pageable.Unpaged(), sortable.Unsorted(), filterable.New(idFilter)))
	if err != nil {
		return Account{}, fmt.Errorf("failed to search accounts: %w", err)
	}
	if accounts.NumberOfElements() != 1 {
		return Account{}, commonErrors.ErrRecordNotFound()
	}

	return accounts.Content()[0], nil
}
