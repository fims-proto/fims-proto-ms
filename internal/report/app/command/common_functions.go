package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/template"

	"github.com/google/uuid"
)

func prepareTables(
	ctx context.Context,
	repo domain.Repository,
	sobId uuid.UUID,
	commands []TableCmd,
) ([]*template.Table, error) {
	return nil, nil
}

func prepareLineItems(
	commands []LineItemCmd,
) ([]*template.Item, error) {
	return nil, nil
}
