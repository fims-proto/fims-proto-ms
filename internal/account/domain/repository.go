package domain

import "context"

type Repository interface {
	AddAccount(ctx context.Context, account *Account) error
	Dataload(ctx context.Context, accounts []*Account) error
	Migrate(ctx context.Context) error
}
