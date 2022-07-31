package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/journal/app/service"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"

	"github/fims-proto/fims-proto-ms/internal/journal/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type PostJournalEntryCmd struct {
	EntryId uuid.UUID
	Poster  uuid.UUID
}

type PostJournalEntryHandler struct {
	repo           domain.Repository
	accountService service.AccountService
}

func NewPostJournalEntryHandler(repo domain.Repository, accountService service.AccountService) PostJournalEntryHandler {
	if repo == nil {
		panic("nil repo")
	}

	if accountService == nil {
		panic("nil account service")
	}

	return PostJournalEntryHandler{
		repo:           repo,
		accountService: accountService,
	}
}

func (h PostJournalEntryHandler) Handle(ctx context.Context, cmd PostJournalEntryCmd) error {
	return h.repo.UpdateJournalEntry(
		ctx,
		cmd.EntryId,
		func(j *journal_entry.JournalEntry) (*journal_entry.JournalEntry, error) {
			if err := j.Post(cmd.Poster); err != nil {
				return nil, err
			}

			if err := h.accountService.PostJournalEntry(ctx, *j); err != nil {
				return nil, errors.Wrap(err, "failed to post journal entry to accounts")
			}

			return j, nil
		},
	)
}
