package account

import (
	"context"

	"github.com/pkg/errors"

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

func (i IntraProcessAdapter) InitializeFirstPeriod(ctx context.Context, sobId uuid.UUID, fiscalYear, number int) error {
	periodId := uuid.New()
	if err := i.accountInterface.CreatePeriodByNumber(ctx, command.CreateCurrentPeriodCmd{
		SobId:      sobId,
		PeriodId:   periodId,
		FiscalYear: fiscalYear,
		Number:     number,
	}); err != nil {
		return errors.Wrap(err, "initializing first period failed")
	}

	if err := i.accountInterface.CreateLedgers(ctx, command.CreateLedgersCmd{
		SobId:    sobId,
		PeriodId: periodId,
	}); err != nil {
		return errors.Wrap(err, "initializing ledgers for first period failed")
	}

	return nil
}
