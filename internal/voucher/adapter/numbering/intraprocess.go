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

func (i IntraProcessAdapter) GenerateIdentifier(ctx context.Context, periodId uuid.UUID, voucherType string) (string, error) {
	return i.numberingInterface.GenerateIdentifierForVoucher(ctx, periodId, voucherType)
}
