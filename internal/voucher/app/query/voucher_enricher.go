package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	userQuery "github/fims-proto/fims-proto-ms/internal/user/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/service"
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

func enrichLineItemAccountNumber(ctx context.Context, service service.AccountService, vouchers []Voucher) ([]Voucher, error) {
	accountSet := make(map[uuid.UUID]void)
	for _, v := range vouchers {
		for _, item := range v.LineItems {
			accountSet[item.AccountId] = empty
		}
	}

	accountConfigs, err := service.ReadAccountsByIds(ctx, toKeySlice[uuid.UUID, void](accountSet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read accounts by Ids")
	}

	for i := range vouchers {
		for j := range vouchers[i].LineItems {
			accountConfig, ok := accountConfigs[vouchers[i].LineItems[j].AccountId]
			if !ok {
				return nil, errors.Errorf("account not found by id: %s", vouchers[i].LineItems[j].AccountId)
			}
			vouchers[i].LineItems[j].AccountNumber = accountConfig.AccountNumber
		}
	}

	return vouchers, nil
}

func enrichUserName(ctx context.Context, service service.UserService, vouchers []Voucher) ([]Voucher, error) {
	userSet := make(map[uuid.UUID]void)
	for _, v := range vouchers {
		if v.Creator.Id != uuid.Nil {
			userSet[v.Creator.Id] = empty
		}
		if v.Reviewer.Id != uuid.Nil {
			userSet[v.Reviewer.Id] = empty
		}
		if v.Auditor.Id != uuid.Nil {
			userSet[v.Auditor.Id] = empty
		}
		if v.Poster.Id != uuid.Nil {
			userSet[v.Poster.Id] = empty
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

	for i := range vouchers {
		if vouchers[i].Creator.Id != uuid.Nil {
			vouchers[i].Creator = convertUser(vouchers[i].Creator, users)
		}
		if vouchers[i].Reviewer.Id != uuid.Nil {
			vouchers[i].Reviewer = convertUser(vouchers[i].Reviewer, users)
		}
		if vouchers[i].Auditor.Id != uuid.Nil {
			vouchers[i].Auditor = convertUser(vouchers[i].Auditor, users)
		}
		if vouchers[i].Poster.Id != uuid.Nil {
			vouchers[i].Poster = convertUser(vouchers[i].Poster, users)
		}
	}

	return vouchers, nil
}

func enrichPeriod(ctx context.Context, service service.AccountService, vouchers []Voucher) ([]Voucher, error) {
	periodSet := make(map[uuid.UUID]void)
	for _, v := range vouchers {
		periodSet[v.Period.PeriodId] = empty
	}

	periods, err := service.ReadPeriodsByIds(ctx, toKeySlice[uuid.UUID, void](periodSet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to read periods by Ids")
	}

	for i := range vouchers {
		period, ok := periods[vouchers[i].Period.PeriodId]
		if !ok {
			return nil, errors.Errorf("period not found by id: %s", vouchers[i].Period.PeriodId)
		}
		vouchers[i].Period = Period{
			PeriodId:    period.Id,
			FiscalYear:  period.FiscalYear,
			Number:      period.PeriodNumber,
			OpeningTime: period.OpeningTime,
			EndingTime:  period.EndingTime,
			IsClosed:    period.IsClosed,
			CreatedAt:   period.CreatedAt,
			UpdatedAt:   period.UpdatedAt,
		}
	}

	return vouchers, nil
}
