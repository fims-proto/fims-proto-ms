package domain

import "context"

type Repository interface {
	AddAccount(ctx context.Context, account *Account) error
	AddAccounts(ctx context.Context, accounts []*Account) error
}
