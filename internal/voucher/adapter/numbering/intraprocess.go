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
