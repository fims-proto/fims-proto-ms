package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"

	"github.com/google/uuid"
)

type AuditJournalCmd struct {
	JournalId uuid.UUID
	Auditor   uuid.UUID
}

type AuditJournalHandler struct {
	repo domain.Repository
}

func NewAuditJournalHandler(repo domain.Repository) AuditJournalHandler {
	if repo == nil {
		panic("nil repo")
	}

	return AuditJournalHandler{repo: repo}
}

func (h AuditJournalHandler) Handle(ctx context.Context, cmd AuditJournalCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.audit(txCtx, cmd)
	})
}

func (h AuditJournalHandler) audit(ctx context.Context, cmd AuditJournalCmd) error {
	return h.repo.UpdateJournalHeader(
		ctx,
		cmd.JournalId,
		func(j *journal.Journal) (*journal.Journal, error) {
			err := j.Audit(cmd.Auditor)
			return j, err
		},
	)
}
