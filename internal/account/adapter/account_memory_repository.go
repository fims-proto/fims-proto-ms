package adapter

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"sync"

	"github.com/pkg/errors"
)

type AccountMemoryRepository struct {
	lock *sync.RWMutex
	data map[string]account.Account
}

func NewAccountMemoryRepository() AccountMemoryRepository {
	return AccountMemoryRepository{
		lock: &sync.RWMutex{},
		data: make(map[string]account.Account),
	}
}

func (r AccountMemoryRepository) ValidateExistence(ctx context.Context, accNumbers []string) error {
	// fake logic for now
	for _, accNumber := range accNumbers {
		if accNumber == "0000" {
			return errors.Errorf("invalid account number %s", accNumber)
		}
	}
	return nil
}
