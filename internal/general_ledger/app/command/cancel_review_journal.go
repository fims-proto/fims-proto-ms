package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
)

type CancelReviewJournalCmd struct {
	JournalId uuid.UUID
	Reviewer  uuid.UUID
}

type CancelReviewJournalHandler struct {
	repo domain.Repository
}

func NewCancelReviewJournalHandler(repo domain.Repository) CancelReviewJournalHandler {
	if repo == nil {
		panic("nil repo")
	}

	return CancelReviewJournalHandler{repo: repo}
}

func (h CancelReviewJournalHandler) Handle(ctx context.Context, cmd CancelReviewJournalCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.cancelReview(txCtx, cmd)
	})
}

func (h CancelReviewJournalHandler) cancelReview(ctx context.Context, cmd CancelReviewJournalCmd) error {
	return h.repo.UpdateJournalHeader(
		ctx,
		cmd.JournalId,
		func(j *journal.Journal) (*journal.Journal, error) {
			err := j.CancelReview(cmd.Reviewer)
			return j, err
		},
	)
}
