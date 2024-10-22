package query

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type ReportReadModel interface {
	SearchReport(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Report], error)
}
