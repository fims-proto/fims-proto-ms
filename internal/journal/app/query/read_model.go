package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type JournalReadModel interface {
	SearchJournalEntries(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[JournalEntry], error)

	JournalEntryById(ctx context.Context, entryId uuid.UUID) (JournalEntry, error)
}
