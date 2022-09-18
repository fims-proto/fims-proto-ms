package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/datav3"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/journal/app/service"
)

type PagingJournalEntriesHandler struct {
	readModel   JournalReadModel
	userService service.UserService
}

func NewPagingJournalEntriesHandler(readModel JournalReadModel, userService service.UserService) PagingJournalEntriesHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if userService == nil {
		panic("nil user service")
	}

	return PagingJournalEntriesHandler{
		readModel:   readModel,
		userService: userService,
	}
}

func (h PagingJournalEntriesHandler) Handle(ctx context.Context, sobId uuid.UUID, pageRequest datav3.PageRequest) (datav3.Page[JournalEntry], error) {
	entriesPage, err := h.readModel.SearchJournalEntries(ctx, sobId, pageRequest)
	if err != nil {
		return nil, err
	}

	journalEntries := entriesPage.Content()

	journalEntries, err = enrichUserName(ctx, h.userService, journalEntries)
	if err != nil {
		return nil, errors.Wrap(err, "failed to enrich user in journal entry")
	}

	return datav3.NewPage(journalEntries, pageRequest, entriesPage.NumberOfElements())
}
