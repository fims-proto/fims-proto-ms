package sob

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"

	sobPort "github/fims-proto/fims-proto-ms/internal/sob/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	sobInterface sobPort.SobInterface
}

func NewIntraProcessAdapter(sobInterface sobPort.SobInterface) IntraProcessAdapter {
	return IntraProcessAdapter{sobInterface: sobInterface}
}

func (s IntraProcessAdapter) ReadById(ctx context.Context, sobId uuid.UUID) (query.Sob, error) {
	return s.sobInterface.ReadById(ctx, sobId)
}
