package command

import (
	"context"
	"fmt"

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
)([]*template.Item, error)
{

}

func prepareReportInnerTemplate(
	ctx context.Context,
	repo domain.Repository,
	sobId uuid.UUID,
	refTemplateId uuid.UUID,
) (*Template, error) {
	// validate templateId
	refTemplate, err := repo.ReadTemplateById(ctx, refTemplateId)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err);
	}
}
