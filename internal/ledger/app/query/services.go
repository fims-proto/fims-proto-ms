package query

import (
	"context"

	"github.com/google/uuid"
	accountQuery "github/fims-proto/fims-proto-ms/internal/account/app/query"
)

type AccountService interface {
	ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]accountQuery.Account, error)
}
