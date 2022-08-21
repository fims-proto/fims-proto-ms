package query

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type JournalReadModel interface {
	JournalEntryById(ctx context.Context, entryId uuid.UUID) (JournalEntry, error)
	PagingJournalEntries(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[JournalEntry], error)
}
