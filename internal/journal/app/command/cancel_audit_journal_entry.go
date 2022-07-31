package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"

	"github/fims-proto/fims-proto-ms/internal/journal/domain"

	"github.com/google/uuid"
)

type CancelAuditJournalEntryCmd struct {
	EntryId uuid.UUID
	Auditor uuid.UUID
}

type CancelAuditJournalEntryHandler struct {
	repo domain.Repository
}

func NewCancelAuditJournalEntryHandler(repo domain.Repository) CancelAuditJournalEntryHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CancelAuditJournalEntryHandler{repo: repo}
}

func (h CancelAuditJournalEntryHandler) Handle(ctx context.Context, cmd CancelAuditJournalEntryCmd) error {
	return h.repo.UpdateJournalEntry(
		ctx,
		cmd.EntryId,
		func(j *journal_entry.JournalEntry) (*journal_entry.JournalEntry, error) {
			err := j.CancelAudit(cmd.Auditor)
			return j, err
		},
	)
}
