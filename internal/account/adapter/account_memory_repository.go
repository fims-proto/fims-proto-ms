package adapter

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
	"sync"

	"github.com/pkg/errors"
)

type AccountMemoryRepository struct {
	lock *sync.RWMutex
	data map[string]domain.Account
}

func NewAccountMemoryRepository() AccountMemoryRepository {
	return AccountMemoryRepository{
		lock: &sync.RWMutex{},
		data: make(map[string]domain.Account),
	}
}

func (r AccountMemoryRepository) ValidateExistence(ctx context.Context, accNumbers []string) error {
	// TODO: fake logic for now
	for _, accNumber := range accNumbers {
		if accNumber == "0000" {
			return errors.Errorf("invalid account number %s", accNumber)
		}
	}
	return nil
}

func (r AccountMemoryRepository) AllAccounts(ctx context.Context) ([]query.Account, error) {
	// TODO: fake logic for now
	panic("not implemented")
}

func (r AccountMemoryRepository) AccountByNumber(ctx context.Context, accountNumber string) (query.Account, error) {
	// TODO: fake logic for now
	if len(accountNumber) != 8 {
		return query.Account{}, errors.New("let's test with 8 length account number")
	}
	return query.Account{
		Number:      accountNumber,
		Title:       "3rd lvl",
		AccountType: "assets",
		SuperiorAccount: &query.Account{
			Number:      accountNumber[:6],
			Title:       "2nd lvl",
			AccountType: "assets",
			SuperiorAccount: &query.Account{
				Number:      accountNumber[:4],
				Title:       "1st lvl",
				AccountType: "assets",
			},
		},
	}, nil
}
