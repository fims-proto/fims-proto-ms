package numbering

import (
	"context"

	"github.com/google/uuid"
	numberingPort "github/fims-proto/fims-proto-ms/internal/numbering/port/private/intraprocess"
)

type IntraProcessAdapter struct {
	numberingInterface numberingPort.NumberingInterface
}

func NewIntraProcessAdapter(numberingInterface numberingPort.NumberingInterface) IntraProcessAdapter {
	return IntraProcessAdapter{numberingInterface: numberingInterface}
}

func (i IntraProcessAdapter) InitializeIdentifierConfigurationForVoucher(ctx context.Context, periodId uuid.UUID) error {
	return i.numberingInterface.CreateIdentifierConfigurationForVoucher(ctx, periodId, "general_journal")
}
