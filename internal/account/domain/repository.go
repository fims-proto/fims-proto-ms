package domain

import "context"

type Repository interface {
	CreateAccount(ctx context.Context, account *Account) error
	DataLoad(ctx context.Context, accounts []*Account) error
	Migrate(ctx context.Context) error
}
