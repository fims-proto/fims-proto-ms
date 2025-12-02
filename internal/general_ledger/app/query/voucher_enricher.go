package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/utils"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github.com/google/uuid"
)

type void struct{}

var empty void

func enrichUserName(ctx context.Context, service service.UserService, vouchers []Voucher) ([]Voucher, error) {
	userSet := make(map[uuid.UUID]void)
	addSetIfNotNil := func(u *User) {
		if u != nil {
			userSet[u.Id] = empty
		}
	}
	for _, v := range vouchers {
		addSetIfNotNil(v.Creator)
		addSetIfNotNil(v.Reviewer)
		addSetIfNotNil(v.Auditor)
		addSetIfNotNil(v.Poster)
	}

	users, err := service.ReadUsersByIds(ctx, utils.MapToKeySlice(userSet))
	if err != nil {
		return nil, fmt.Errorf("failed to read users: %w", err)
	}

	enrichUser := func(u *User) *User {
		if u != nil {
			return &User{
				Id:     u.Id,
				Traits: users[u.Id].Traits,
			}
		}
		return nil
	}

	for i := range vouchers {
		vouchers[i].Creator = enrichUser(vouchers[i].Creator)
		vouchers[i].Reviewer = enrichUser(vouchers[i].Reviewer)
		vouchers[i].Auditor = enrichUser(vouchers[i].Auditor)
		vouchers[i].Poster = enrichUser(vouchers[i].Poster)
	}

	return vouchers, nil
}
