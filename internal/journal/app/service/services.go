package service

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"

	accountQuery "github/fims-proto/fims-proto-ms/internal/account/app/query"
	userQuery "github/fims-proto/fims-proto-ms/internal/user/app/query"

	"github.com/google/uuid"
)

type AccountService interface {
	ValidateExistenceAndGetId(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]uuid.UUID, error)
	ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]accountQuery.Account, error)
	ReadPeriodByTime(ctx context.Context, sobId uuid.UUID, transactionTime time.Time) (accountQuery.Period, error)
	ReadPeriodsByIds(ctx context.Context, periodIds []uuid.UUID) (map[uuid.UUID]accountQuery.Period, error)
	PostJournalEntry(ctx context.Context, entry journal_entry.JournalEntry) error
}

type UserService interface {
	ReadUsersByIds(ctx context.Context, userIds []uuid.UUID) (map[uuid.UUID]userQuery.User, error)
}

type NumberingService interface {
	GenerateIdentifier(ctx context.Context, periodId uuid.UUID, journalType string) (string, error)
}
