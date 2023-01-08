package account

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/account/app/command"

	accountPort "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	accountInterface accountPort.AccountInterface
}

func NewIntraProcessAdapter(accountInterface accountPort.AccountInterface) IntraProcessAdapter {
	return IntraProcessAdapter{accountInterface: accountInterface}
}

func (i IntraProcessAdapter) InitializeAccounts(ctx context.Context, sobId uuid.UUID) error {
	return i.accountInterface.InitializeAccounts(ctx, sobId)
}

func (i IntraProcessAdapter) InitializeFirstPeriod(ctx context.Context, sobId uuid.UUID, financialYear, number int) error {
	startDateOfMonth := time.Date(financialYear, time.Month(number), 1, 0, 0, 0, 0, time.UTC)

	return i.accountInterface.CreatePeriod(ctx, command.CreatePeriodCmd{
		SobId:         sobId,
		PeriodId:      uuid.New(),
		FinancialYear: financialYear,
		Number:        number,
		OpeningTime:   startDateOfMonth,
	})
}
