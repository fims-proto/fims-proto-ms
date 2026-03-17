package numbering

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/numbering/app/command"

	numberingPort "github/fims-proto/fims-proto-ms/internal/numbering/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	numberingInterface numberingPort.NumberingInterface
}

func NewIntraProcessAdapter(numberingInterface numberingPort.NumberingInterface) IntraProcessAdapter {
	return IntraProcessAdapter{numberingInterface: numberingInterface}
}

func (i IntraProcessAdapter) GenerateIdentifier(ctx context.Context, periodId uuid.UUID) (string, error) {
	cmd := command.GenerateNextIdentifierCmd{
		IdentifierId:         uuid.New(),
		TargetBusinessObject: "journal",
		ObjectsToMatch: map[string]string{
			"period_id": periodId.String(),
		},
	}

	return i.numberingInterface.GenerateIdentifier(ctx, cmd)
}

func (i IntraProcessAdapter) CreateIdentifierConfigurationForJournal(ctx context.Context, periodId uuid.UUID) error {
	cmd := command.CreateIdentifierConfigurationCmd{
		IdentifierConfigurationId: uuid.New(),
		TargetBusinessObject:      "journal",
		PropertyMatchers: []struct{ Name, Value string }{
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
