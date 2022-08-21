package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"

	"github/fims-proto/fims-proto-ms/internal/journal/domain"

	"github.com/google/uuid"
)

type ReviewJournalEntryCmd struct {
	EntryId  uuid.UUID
	Reviewer uuid.UUID
}

type ReviewJournalEntryHandler struct {
	repo domain.Repository
}

func NewReviewJournalEntryHandler(repo domain.Repository) ReviewJournalEntryHandler {
	if repo == nil {
		panic("nil repo")
	}

	return ReviewJournalEntryHandler{repo: repo}
}

func (h ReviewJournalEntryHandler) Handle(ctx context.Context, cmd ReviewJournalEntryCmd) error {
	return h.repo.UpdateJournalEntry(
		ctx,
		cmd.EntryId,
		func(j *journal_entry.JournalEntry) (*journal_entry.JournalEntry, error) {
			err := j.Review(cmd.Reviewer)
			return j, err
		},
	)
}
