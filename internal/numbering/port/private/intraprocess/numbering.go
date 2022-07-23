package intraprocess

import (
	"context"

	"github.com/google/uuid"
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

func (i NumberingInterface) CreateIdentifierConfigurationForVoucher(ctx context.Context, periodId uuid.UUID, voucherType string) error {
	cmd := command.CreateIdentifierConfigurationCmd{
		Id:                   uuid.New(),
		TargetBusinessObject: "voucher",
		PropertyMatchers: []struct{ Name, Value string }{
			{
				Name:  "voucher_type",
				Value: voucherType,
			},
			{
				Name:  "period_id",
				Value: periodId.String(),
			},
		},
		Prefix: "记 ",
		Suffix: " 号",
	}
	return i.app.Commands.CreateIdentifierConfiguration.Handle(ctx, cmd)
}

func (i NumberingInterface) GenerateIdentifierForVoucher(ctx context.Context, periodId uuid.UUID, voucherType string) (string, error) {
	createdId := uuid.New()
	cmd := command.GenerateNextIdentifierCmd{
		IdentifierId:         createdId,
		TargetBusinessObject: "voucher",
		ObjectsToMatch: map[string]string{
			"voucher_type": voucherType,
			"period_id":    periodId.String(),
		},
	}
	if err := i.app.Commands.GenerateNextIdentifier.Handle(ctx, cmd); err != nil {
		return "", errors.Wrap(err, "failed to generate identifier")
	}

	identifier, err := i.app.Queries.IdentifierById.Handle(ctx, createdId)
	if err != nil {
		return "", errors.Wrap(err, "failed to read identifier")
	}

	return identifier.Identifier, nil
}
