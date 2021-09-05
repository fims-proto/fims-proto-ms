package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
)

type MigrationHanlder struct {
	repo domain.Repository
}

func NewMigrationHanlder(repo domain.Repository) MigrationHanlder {
	if repo == nil {
		panic("nil repo")
	}
	return MigrationHanlder{repo: repo}
}

func (h MigrationHanlder) Handle(ctx context.Context) error {
	return h.repo.Migrate(ctx)
}
