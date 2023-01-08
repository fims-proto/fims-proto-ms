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

func (i IntraProcessAdapter) InitializeIdentifierConfigurationForJournal(ctx context.Context, periodId uuid.UUID) error {
	cmd := command.CreateIdentifierConfigurationCmd{
		IdentifierConfigurationId: uuid.New(),
		TargetBusinessObject:      "journal",
		PropertyMatchers: []struct{ Name, Value string }{
			{
				Name:  "journal_type",
				Value: "general_journal",
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
