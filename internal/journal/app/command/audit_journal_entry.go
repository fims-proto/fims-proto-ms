package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/journal/domain"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"

	"github.com/google/uuid"
)

type AuditJournalEntryCmd struct {
	EntryId uuid.UUID
	Auditor uuid.UUID
}

type AuditJournalEntryHandler struct {
	repo domain.Repository
}

func NewAuditJournalEntryHandler(repo domain.Repository) AuditJournalEntryHandler {
	if repo == nil {
		panic("nil repo")
	}
	return AuditJournalEntryHandler{repo: repo}
}

func (h AuditJournalEntryHandler) Handle(ctx context.Context, cmd AuditJournalEntryCmd) error {
	return h.repo.UpdateJournalEntry(
		ctx,
		cmd.EntryId,
		func(j *journal_entry.JournalEntry) (*journal_entry.JournalEntry, error) {
			err := j.Audit(cmd.Auditor)
			return j, err
		},
	)
}
