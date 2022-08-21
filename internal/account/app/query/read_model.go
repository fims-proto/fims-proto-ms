package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type AccountReadModel interface {
	AllAccountConfigurations(ctx context.Context, sobId uuid.UUID) ([]AccountConfiguration, error)
	PagingAccountConfigurations(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[AccountConfiguration], error)
	AccountConfigurationsByIds(ctx context.Context, accountIds []uuid.UUID) ([]AccountConfiguration, error)
	AccountConfigurationsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]AccountConfiguration, error)
	SuperiorAccountConfigurations(ctx context.Context, accountId uuid.UUID) ([]AccountConfiguration, error)

	AccountsInPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID) ([]Account, error)

	PeriodById(ctx context.Context, periodId uuid.UUID) (Period, error)
	PeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (Period, error)
	PeriodsByIds(ctx context.Context, periodIds []uuid.UUID) ([]Period, error)
}
