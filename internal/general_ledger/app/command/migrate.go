package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
)

type MigrationHandler struct {
	repo domain.Repository
}

func NewMigrationHandler(repo domain.Repository) MigrationHandler {
	if repo == nil {
		panic("nil repo")
	}
	return MigrationHandler{repo: repo}
}

func (h MigrationHandler) Handle(ctx context.Context) error {
	return h.repo.Migrate(ctx)
}
