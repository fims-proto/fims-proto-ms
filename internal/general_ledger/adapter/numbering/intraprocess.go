package numbering

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/numbering/app/command"

	"github.com/google/uuid"
	numberingPort "github/fims-proto/fims-proto-ms/internal/numbering/port/private/intraprocess"
)

type IntraProcessAdapter struct {
	numberingInterface numberingPort.NumberingInterface
}

func NewIntraProcessAdapter(numberingInterface numberingPort.NumberingInterface) IntraProcessAdapter {
	return IntraProcessAdapter{numberingInterface: numberingInterface}
}

func (i IntraProcessAdapter) GenerateIdentifier(ctx context.Context, periodId uuid.UUID, voucherType string) (string, error) {
	cmd := command.GenerateNextIdentifierCmd{
		IdentifierId:         uuid.New(),
		TargetBusinessObject: "voucher",
		ObjectsToMatch: map[string]string{
			"voucher_type": voucherType,
			"period_id":    periodId.String(),
		},
	}

	return i.numberingInterface.GenerateIdentifier(ctx, cmd)
}

func (i IntraProcessAdapter) CreateIdentifierConfigurationForVoucher(ctx context.Context, periodId uuid.UUID) error {
	cmd := command.CreateIdentifierConfigurationCmd{
		IdentifierConfigurationId: uuid.New(),
		TargetBusinessObject:      "voucher",
		PropertyMatchers: []struct{ Name, Value string }{
			{
				Name:  "voucher_type",
				Value: "general_voucher",
			},
			{
				Name:  "period_id",
				Value: periodId.String(),
			},
		},
		Prefix: "记 ",
		Suffix: " 号",
	}
	return i.numberingInterface.CreateIdentifierConfiguration(ctx, cmd)
}
