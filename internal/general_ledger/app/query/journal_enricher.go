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

func enrichUserName(ctx context.Context, service service.UserService, journals []Journal) ([]Journal, error) {
	userSet := make(map[uuid.UUID]void)
	addSetIfNotNil := func(u *User) {
		if u != nil {
			userSet[u.Id] = empty
		}
	}
	for _, j := range journals {
		addSetIfNotNil(j.Creator)
		addSetIfNotNil(j.Reviewer)
		addSetIfNotNil(j.Auditor)
		addSetIfNotNil(j.Poster)
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

	for i := range journals {
		journals[i].Creator = enrichUser(journals[i].Creator)
		journals[i].Reviewer = enrichUser(journals[i].Reviewer)
		journals[i].Auditor = enrichUser(journals[i].Auditor)
		journals[i].Poster = enrichUser(journals[i].Poster)
	}

	return journals, nil
}
