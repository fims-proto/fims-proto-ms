package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/utils"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	userQuery "github/fims-proto/fims-proto-ms/internal/user/app/query"
)

type void struct{}

var empty void

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

	users, err := service.ReadUsersByIds(ctx, utils.MapToKeySlice(userSet))
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
