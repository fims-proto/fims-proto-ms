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
	sobService  service.SobService
	userService service.UserService
}

func NewPagingJournalsHandler(readModel GeneralLedgerReadModel, sobService service.SobService, userService service.UserService) PagingJournalsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	if userService == nil {
		panic("nil user service")
	}

	return PagingJournalsHandler{
		readModel:   readModel,
		sobService:  sobService,
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

	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return nil, fmt.Errorf("failed to read sob: %w", err)
	}

	journals = enrichJournalAccountNumbers(sob.AccountsCodeLength, journals)

	return data.NewPage(journals, pageRequest, entriesPage.NumberOfElements())
}
