package query

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type SobReadModel interface {
	PagingSobs(ctx context.Context, pageable data.Pageable) (data.Page[Sob], error)
	SobById(ctx context.Context, sobId uuid.UUID) (Sob, error)
}
