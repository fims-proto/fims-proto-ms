package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github.com/google/uuid"
)

type PagingJournalsHandler struct {
	readModel   GeneralLedgerReadModel
	userService service.UserService
}

func NewPagingJournalsHandler(readModel GeneralLedgerReadModel, userService service.UserService) PagingJournalsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if userService == nil {
		panic("nil user service")
	}

	return PagingJournalsHandler{
		readModel:   readModel,
		userService: userService,
	}
}

func (h PagingJournalsHandler) Handle(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Journal], error) {
	entriesPage, err := h.readModel.SearchJournals(ctx, sobId, pageRequest)
	if err != nil {
		return nil, err
	}

	journals := entriesPage.Content()

	journals, err = enrichUserName(ctx, h.userService, journals)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich user in journal: %w", err)
	}

	return data.NewPage(journals, pageRequest, entriesPage.NumberOfElements())
}
