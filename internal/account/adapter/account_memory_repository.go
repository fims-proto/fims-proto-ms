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

func (r AccountMemoryRepository) AllAccounts(ctx context.Context) ([]query.Account, error) {
	panic("not implemented")
}

func (r AccountMemoryRepository) AccountByNumber(ctx context.Context, accountNumber string) (query.Account, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	result, err := r.readAccountWithSuperiorAccount(accountNumber)
	if err != nil {
		return query.Account{}, errors.Wrapf(err, "failed to read account by number %s", accountNumber)
	}
	return result, nil
}

func (r AccountMemoryRepository) AddAccount(ctx context.Context, account *domain.Account) error {
	panic("not implemented")
}

func (r AccountMemoryRepository) Dataload(ctx context.Context, accounts []*domain.Account) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	// clear map
	for i := range r.data {
		delete(r.data, i)
	}

	for _, account := range accounts {
		r.data[account.Number()] = *account
	}
	return nil
}

func (r AccountMemoryRepository) readAccountWithSuperiorAccount(accountNumber string) (query.Account, error) {
	account, ok := r.data[accountNumber]
	if !ok {
		return query.Account{}, errors.Errorf("account number %s does not exist", accountNumber)
	}
	result := mapFromDomainAccount(account)
	if account.SuperiorNumber() == "" {
		return result, nil
	}
	superiorAccount, err := r.readAccountWithSuperiorAccount(account.SuperiorNumber())
	if err != nil {
		return query.Account{}, err
	}
	result.SuperiorAccount = &superiorAccount
	return result, nil
}

func mapFromDomainAccount(account domain.Account) query.Account {
	return query.Account{
		Number:          account.Number(),
		Title:           account.Title(),
		AccountType:     account.Type().String(),
		SuperiorAccount: nil,
	}
}
