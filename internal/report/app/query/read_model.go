package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type ReportReadModel interface {
	SearchReport(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Report], error)
}
