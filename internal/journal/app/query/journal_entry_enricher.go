package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/journal/app/service"
	userQuery "github/fims-proto/fims-proto-ms/internal/user/app/query"
)

type void struct{}

var empty void

func toKeySlice[K comparable, V interface{}](set map[K]V) []K {
	keys := make([]K, len(set))
	i := 0
	for k := range set {
		keys[i] = k
		i++
	}
	return keys
}

func enrichLineItemAccountNumber(ctx context.Context, service service.AccountService, entries []JournalEntry) ([]JournalEntry, error) {
	accountSet := make(map[uuid.UUID]void)
	for _, entry := range entries {
		for _, item := range entry.LineItems {
			accountSet[item.AccountId] = empty
		}
	}

	accountConfigs, err := service.ReadAccountsByIds(ctx, toKeySlice[uuid.UUID, void](accountSet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read accounts by Ids")
	}

	for i := range entries {
		for j := range entries[i].LineItems {
			accountConfig, ok := accountConfigs[entries[i].LineItems[j].AccountId]
			if !ok {
				return nil, errors.Errorf("account not found by id: %s", entries[i].LineItems[j].AccountId)
			}
			entries[i].LineItems[j].AccountNumber = accountConfig.AccountNumber
		}
	}

	return entries, nil
}

func enrichUserName(ctx context.Context, service service.UserService, entries []JournalEntry) ([]JournalEntry, error) {
	userSet := make(map[uuid.UUID]void)
	for _, entry := range entries {
		if entry.Creator.Id != uuid.Nil {
			userSet[entry.Creator.Id] = empty
		}
		if entry.Reviewer.Id != uuid.Nil {
			userSet[entry.Reviewer.Id] = empty
		}
		if entry.Auditor.Id != uuid.Nil {
			userSet[entry.Auditor.Id] = empty
		}
	}

	users, err := service.ReadUsersByIds(ctx, toKeySlice[uuid.UUID, void](userSet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read users by Ids")
	}

	convertUser := func(user User, users map[uuid.UUID]userQuery.User) User {
		return User{
			Id:     user.Id,
			Traits: users[user.Id].Traits,
		}
	}

	for i := range entries {
		if entries[i].Creator.Id != uuid.Nil {
			entries[i].Creator = convertUser(entries[i].Creator, users)
		}
		if entries[i].Reviewer.Id != uuid.Nil {
			entries[i].Reviewer = convertUser(entries[i].Reviewer, users)
		}
		if entries[i].Auditor.Id != uuid.Nil {
			entries[i].Auditor = convertUser(entries[i].Auditor, users)
		}
	}

	return entries, nil
}

func enrichPeriod(ctx context.Context, service service.AccountService, entries []JournalEntry) ([]JournalEntry, error) {
	periodSet := make(map[uuid.UUID]void)
	for _, entry := range entries {
		periodSet[entry.Period.PeriodId] = empty
	}

	periods, err := service.ReadPeriodsByIds(ctx, toKeySlice[uuid.UUID, void](periodSet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read periods by Ids")
	}

	for i := range entries {
		period, ok := periods[entries[i].Period.PeriodId]
		if !ok {
			return nil, errors.Errorf("period not found by id: %s", entries[i].Period.PeriodId)
		}
		entries[i].Period = Period{
			PeriodId:      period.Id,
			FinancialYear: period.FiscalYear,
			Number:        period.PeriodNumber,
			OpeningTime:   period.OpeningTime,
			EndingTime:    period.EndingTime,
			IsClosed:      period.IsClosed,
			CreatedAt:     period.CreatedAt,
			UpdatedAt:     period.UpdatedAt,
		}
	}

	return entries, nil
}
