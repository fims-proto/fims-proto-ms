package account

import (
	"context"
	accountPort "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
)

type IntraprocessService struct {
	accInterface accountPort.AccountInterface
}

func NewIntraprocessService(accInterface accountPort.AccountInterface) IntraprocessService {
	return IntraprocessService{accInterface: accInterface}
}

func (s IntraprocessService) ValidateExistence(ctx context.Context, accNumbers []string) error {
	return s.accInterface.ValidateExistence(ctx, accNumbers)
}
