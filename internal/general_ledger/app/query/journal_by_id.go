package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github.com/google/uuid"
)

type JournalByIdHandler struct {
	readModel        GeneralLedgerReadModel
	sobService       service.SobService
	userService      service.UserService
	dimensionService service.DimensionService
}

func NewJournalByIdHandler(
	readModel GeneralLedgerReadModel,
	sobService service.SobService,
	userService service.UserService,
	dimensionService service.DimensionService,
) JournalByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	if userService == nil {
		panic("nil user service")
	}

	if dimensionService == nil {
		panic("nil dimension service")
	}

	return JournalByIdHandler{
		readModel:        readModel,
		sobService:       sobService,
		userService:      userService,
		dimensionService: dimensionService,
	}
}

func (h JournalByIdHandler) Handle(ctx context.Context, journalId uuid.UUID) (Journal, error) {
	v, err := h.readModel.JournalById(ctx, journalId)
	if err != nil {
		return Journal{}, fmt.Errorf("failed to read journal: %w", err)
	}

	singletonList, err := enrichUserName(ctx, h.userService, []Journal{v})
	if err != nil {
		return Journal{}, fmt.Errorf("failed to enrich user in journal: %w", err)
	}

	journal := singletonList[0]

	journal.JournalLines, err = enrichJournalLineDimensionOptions(ctx, h.dimensionService, journal.JournalLines)
	if err != nil {
		return Journal{}, fmt.Errorf("failed to enrich dimension options in journal lines: %w", err)
	}

	for i, line := range journal.JournalLines {
		journal.JournalLines[i].Account, err = enrichAccountDimensionCategories(ctx, h.dimensionService, line.Account)
		if err != nil {
			return Journal{}, fmt.Errorf("failed to enrich dimension categories in journal line account: %w", err)
		}
	}

	sob, err := h.sobService.ReadById(ctx, journal.SobId)
	if err != nil {
		return Journal{}, fmt.Errorf("failed to read sob: %w", err)
	}

	journal = enrichJournalAccountNumbers(sob.AccountsCodeLength, []Journal{journal})[0]

	return journal, nil
}
