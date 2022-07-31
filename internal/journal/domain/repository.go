package domain

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"

	"github.com/google/uuid"
)

type Repository interface {
	CreateJournalEntry(ctx context.Context, d *journal_entry.JournalEntry) error
	UpdateJournalEntry(
		ctx context.Context,
		entryId uuid.UUID,
		updateFn func(d *journal_entry.JournalEntry) (*journal_entry.JournalEntry, error),
	) error

	Migrate(ctx context.Context) error
}
