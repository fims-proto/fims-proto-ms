package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/journal/app/service"
)

type JournalEntryByIdHandler struct {
	readModel      JournalReadModel
	accountService service.AccountService
	userService    service.UserService
}

func NewJournalEntryByIdHandler(
	readModel JournalReadModel,
	accountService service.AccountService,
	userService service.UserService,
) JournalEntryByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if accountService == nil {
		panic("nil account service")
	}

	if userService == nil {
		panic("nil user service")
	}

	return JournalEntryByIdHandler{
		readModel:      readModel,
		accountService: accountService,
		userService:    userService,
	}
}

func (h JournalEntryByIdHandler) Handle(ctx context.Context, entryId uuid.UUID) (JournalEntry, error) {
	journalEntry, err := h.readModel.JournalEntryById(ctx, entryId)
	if err != nil {
		return JournalEntry{}, errors.Wrap(err, "failed to read journal entry")
	}

	singletonList, err := enrichLineItemAccountNumber(ctx, h.accountService, []JournalEntry{journalEntry})
	if err != nil {
		return JournalEntry{}, errors.Wrap(err, "failed to enrich account number in journal entry")
	}

	singletonList, err = enrichUserName(ctx, h.userService, singletonList)
	if err != nil {
		return JournalEntry{}, errors.Wrap(err, "failed to enrich user in journal entry")
	}

	singletonList, err = enrichPeriod(ctx, h.accountService, singletonList)
	if err != nil {
		return JournalEntry{}, errors.Wrap(err, "failed to enrich period in journal entry")
	}

	return singletonList[0], nil
}
