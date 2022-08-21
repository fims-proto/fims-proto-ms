package intraprocess

import (
	"context"

	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/numbering/app"
	"github/fims-proto/fims-proto-ms/internal/numbering/app/command"
)

type NumberingInterface struct {
	app *app.Application
}

func NewNumberingInterface(app *app.Application) NumberingInterface {
	return NumberingInterface{app: app}
}

func (i NumberingInterface) CreateIdentifierConfiguration(ctx context.Context, cmd command.CreateIdentifierConfigurationCmd) error {
	return i.app.Commands.CreateIdentifierConfiguration.Handle(ctx, cmd)
}

func (i NumberingInterface) GenerateIdentifier(ctx context.Context, cmd command.GenerateNextIdentifierCmd) (string, error) {
	if err := i.app.Commands.GenerateNextIdentifier.Handle(ctx, cmd); err != nil {
		return "", errors.Wrap(err, "failed to generate identifier")
	}

	identifier, err := i.app.Queries.IdentifierById.Handle(ctx, cmd.IdentifierId)
	if err != nil {
		return "", errors.Wrap(err, "failed to read identifier")
	}

	return identifier.Identifier, nil
}
