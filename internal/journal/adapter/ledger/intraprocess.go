package ledger

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	ledgerQuery "github/fims-proto/fims-proto-ms/internal/ledger/app/query"

	ledgerCommand "github/fims-proto/fims-proto-ms/internal/ledger/app/command"
	ledgerPort "github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	ledgerInterface ledgerPort.LedgerInterface
}

func NewIntraProcessAdapter(ledgerInterface ledgerPort.LedgerInterface) IntraProcessAdapter {
	return IntraProcessAdapter{ledgerInterface: ledgerInterface}
}

func (s IntraProcessAdapter) PostVoucher(ctx context.Context, voucher domain.Voucher) error {
	accountAmountCommands := make(map[uuid.UUID]ledgerCommand.PostLedgersAmountCmd)
	for _, item := range voucher.LineItems() {
		accountAmountCommands[item.AccountId()] = ledgerCommand.PostLedgersAmountCmd{
			Debit:  item.Debit(),
			Credit: item.Credit(),
		}
	}

	command := ledgerCommand.PostLedgersCmd{
		Accounts: accountAmountCommands,
		PeriodId: voucher.PeriodId(),
	}

	return s.ledgerInterface.PostLedgers(ctx, command)
}

func (s IntraProcessAdapter) ReadPeriodByTime(ctx context.Context, sobId uuid.UUID, transactionTime time.Time) (ledgerQuery.Period, error) {
	return s.ledgerInterface.ReadPeriodByTime(ctx, sobId, transactionTime)
}

func (s IntraProcessAdapter) ReadPeriodsByIds(ctx context.Context, periodIds []uuid.UUID) (map[uuid.UUID]ledgerQuery.Period, error) {
	return s.ledgerInterface.ReadPeriodsByIds(ctx, periodIds)
}
