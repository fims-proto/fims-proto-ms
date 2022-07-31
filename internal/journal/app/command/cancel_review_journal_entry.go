package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"

	"github/fims-proto/fims-proto-ms/internal/journal/domain"

	"github.com/google/uuid"
)

type CancelReviewJournalEntryCmd struct {
	EntryId  uuid.UUID
	Reviewer uuid.UUID
}

type CancelReviewJournalEntryHandler struct {
	repo domain.Repository
}

func NewCancelReviewJournalEntryHandler(repo domain.Repository) CancelReviewJournalEntryHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CancelReviewJournalEntryHandler{repo: repo}
}

func (h CancelReviewJournalEntryHandler) Handle(ctx context.Context, cmd CancelReviewJournalEntryCmd) error {
	return h.repo.UpdateJournalEntry(
		ctx,
		cmd.EntryId,
		func(j *journal_entry.JournalEntry) (*journal_entry.JournalEntry, error) {
			err := j.CancelReview(cmd.Reviewer)
			return j, err
		},
	)
}
