package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
)

type ReviewJournalCmd struct {
	JournalId uuid.UUID
	Reviewer  uuid.UUID
}

type ReviewJournalHandler struct {
	repo domain.Repository
}

func NewReviewJournalHandler(repo domain.Repository) ReviewJournalHandler {
	if repo == nil {
		panic("nil repo")
	}

	return ReviewJournalHandler{repo: repo}
}

func (h ReviewJournalHandler) Handle(ctx context.Context, cmd ReviewJournalCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.review(txCtx, cmd)
	})
}

func (h ReviewJournalHandler) review(ctx context.Context, cmd ReviewJournalCmd) error {
	return h.repo.UpdateJournalHeader(
		ctx,
		cmd.JournalId,
		func(j *journal.Journal) (*journal.Journal, error) {
			err := j.Review(cmd.Reviewer)
			return j, err
		},
	)
}
