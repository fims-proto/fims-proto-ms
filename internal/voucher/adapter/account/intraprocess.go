package account

import (
	"context"
	accountport "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
)

type IntraprocessService struct {
	accInterface accountport.AccountInterface
}

func NewIntraprocessService(accInterface accountport.AccountInterface) IntraprocessService {
	return IntraprocessService{accInterface: accInterface}
}

func (s IntraprocessService) ValidateExistence(ctx context.Context, accNumbers []string) error {
	return s.accInterface.ValidateExistence(ctx, accNumbers)
}
