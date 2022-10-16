package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type SobReadModel interface {
	SearchSobs(ctx context.Context, pageRequest data.PageRequest) (data.Page[Sob], error)

	SobById(ctx context.Context, sobId uuid.UUID) (Sob, error)
}
