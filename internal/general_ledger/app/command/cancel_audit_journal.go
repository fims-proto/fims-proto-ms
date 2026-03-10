package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
)

type CancelAuditJournalCmd struct {
	JournalId uuid.UUID
	Auditor   uuid.UUID
}

type CancelAuditJournalHandler struct {
	repo domain.Repository
}

func NewCancelAuditJournalHandler(repo domain.Repository) CancelAuditJournalHandler {
	if repo == nil {
		panic("nil repo")
	}

	return CancelAuditJournalHandler{repo: repo}
}

func (h CancelAuditJournalHandler) Handle(ctx context.Context, cmd CancelAuditJournalCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.cancelAudit(txCtx, cmd)
	})
}

func (h CancelAuditJournalHandler) cancelAudit(ctx context.Context, cmd CancelAuditJournalCmd) error {
	return h.repo.UpdateJournalHeader(
		ctx,
		cmd.JournalId,
		func(j *journal.Journal) (*journal.Journal, error) {
			err := j.CancelAudit(cmd.Auditor)
			return j, err
		},
	)
}
