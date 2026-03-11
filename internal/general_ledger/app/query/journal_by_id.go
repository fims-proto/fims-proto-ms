package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github.com/google/uuid"
)

type JournalByIdHandler struct {
	readModel   GeneralLedgerReadModel
	userService service.UserService
}

func NewJournalByIdHandler(
	readModel GeneralLedgerReadModel,
	userService service.UserService,
) JournalByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if userService == nil {
		panic("nil user service")
	}

	return JournalByIdHandler{
		readModel:   readModel,
		userService: userService,
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

	return singletonList[0], nil
}
