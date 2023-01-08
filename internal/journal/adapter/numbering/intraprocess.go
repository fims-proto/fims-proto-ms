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

func (i IntraProcessAdapter) GenerateIdentifier(ctx context.Context, periodId uuid.UUID, journalType string) (string, error) {
	cmd := command.GenerateNextIdentifierCmd{
		IdentifierId:         uuid.New(),
		TargetBusinessObject: "journal",
		ObjectsToMatch: map[string]string{
			"journal_type": journalType,
			"period_id":    periodId.String(),
		},
	}

	return i.numberingInterface.GenerateIdentifier(ctx, cmd)
}
