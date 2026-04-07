package query

import (
	"context"
	"encoding/json"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github.com/google/uuid"
)

type void struct{}

var empty void

// systemUserTraits is returned for SYSTEM_USER — no service call required.
var systemUserTraits = json.RawMessage(`{"email":"system@fims.internal"}`)

// isSystemUserStub: non-nil *User with Id == uuid.Nil → system actor.
func isSystemUserStub(u *User) bool {
	return u != nil && u.Id == uuid.Nil
}

// enrichUserName enriches journals with full user data (name, traits) by looking up users in the user service
func enrichUserName(ctx context.Context, service service.UserService, journals []Journal) ([]Journal, error) {
	// Collect only real user UUIDs (skip nil = unset, skip system stubs)
	userSet := make(map[uuid.UUID]void)
	for _, j := range journals {
		for _, u := range []*User{j.Creator, j.Reviewer, j.Auditor, j.Poster} {
			if u != nil && u.Id != uuid.Nil {
				userSet[u.Id] = empty
			}
		}
	}

	users, err := service.ReadUsersByIds(ctx, utils.MapToKeySlice(userSet))
	if err != nil {
		return nil, fmt.Errorf("failed to read users: %w", err)
	}

	enrichUser := func(u *User) *User {
		if u == nil {
			return nil
		}
		if isSystemUserStub(u) {
			return &User{Id: uuid.Nil, Traits: systemUserTraits}
		}
		return &User{Id: u.Id, Traits: users[u.Id].Traits}
	}

	for i := range journals {
		journals[i].Creator = enrichUser(journals[i].Creator)
		journals[i].Reviewer = enrichUser(journals[i].Reviewer)
		journals[i].Auditor = enrichUser(journals[i].Auditor)
		journals[i].Poster = enrichUser(journals[i].Poster)
	}

	return journals, nil
}

func enrichAccountDimensionCategories(ctx context.Context, dimensionService service.DimensionService, account Account) (Account, error) {
	if len(account.DimensionCategoryIds) == 0 {
		account.DimensionCategories = []DimensionCategory{}
		return account, nil
	}

	categoriesMap, err := dimensionService.FetchCategoriesByIds(ctx, account.DimensionCategoryIds)
	if err != nil {
		return account, fmt.Errorf("failed to fetch dimension categories: %w", err)
	}

	categories := make([]DimensionCategory, 0, len(account.DimensionCategoryIds))
	for _, id := range account.DimensionCategoryIds {
		if cat, ok := categoriesMap[id]; ok {
			categories = append(categories, DimensionCategory{Id: cat.Id, Name: cat.Name})
		}
	}
	account.DimensionCategories = categories

	return account, nil
}

func enrichJournalLineDimensionOptions(ctx context.Context, dimensionService service.DimensionService, lines []JournalLine) ([]JournalLine, error) {
	// collect all option IDs across all lines
	optionIdSet := make(map[uuid.UUID]void)
	for _, line := range lines {
		for _, id := range line.DimensionOptionIds {
			optionIdSet[id] = empty
		}
	}

	if len(optionIdSet) == 0 {
		for i := range lines {
			lines[i].DimensionOptions = []DimensionOption{}
		}
		return lines, nil
	}

	optionsMap, err := dimensionService.FetchOptionsByIds(ctx, utils.MapToKeySlice(optionIdSet))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dimension options: %w", err)
	}

	for i, line := range lines {
		options := make([]DimensionOption, 0, len(line.DimensionOptionIds))
		for _, id := range line.DimensionOptionIds {
			if opt, ok := optionsMap[id]; ok {
				options = append(options, DimensionOption{
					Id:         opt.Id,
					Name:       opt.Name,
					CategoryId: opt.CategoryId,
					Category:   DimensionCategory{Id: opt.Category.Id, Name: opt.Category.Name},
				})
			}
		}
		lines[i].DimensionOptions = options
	}

	return lines, nil
}
