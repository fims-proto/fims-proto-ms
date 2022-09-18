package query

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/datav3"

	"github.com/google/uuid"
)

type JournalReadModel interface {
	SearchJournalEntries(ctx context.Context, sobId uuid.UUID, pageRequest datav3.PageRequest) (datav3.Page[JournalEntry], error)

	JournalEntryById(ctx context.Context, entryId uuid.UUID) (JournalEntry, error)
}
