package ledger

import (
	"context"
	"time"

	ledgerQuery "github/fims-proto/fims-proto-ms/internal/ledger/app/query"

	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"

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

func (s IntraProcessAdapter) PostVoucher(ctx context.Context, voucher query.Voucher) error {
	// posting id is same in one post
	postingId := uuid.New()
	var commands []ledgerCommand.AppendLedgerLogCmd
	for _, item := range voucher.LineItems {
		commands = append(commands, ledgerCommand.AppendLedgerLogCmd{
			PostingId:       postingId,
			VoucherId:       voucher.Id,
			AccountId:       item.AccountId,
			TransactionTime: voucher.TransactionTime,
			Debit:           item.Debit,
			Credit:          item.Credit,
		})
	}
	return s.ledgerInterface.AppendLedgerLogs(ctx, commands)
}

func (s IntraProcessAdapter) ReadPeriodByTime(ctx context.Context, sobId uuid.UUID, transactionTime time.Time) (ledgerQuery.Period, error) {
	return s.ledgerInterface.ReadPeriodByTime(ctx, sobId, transactionTime)
}

func (s IntraProcessAdapter) ReadPeriodsByIds(ctx context.Context, periodIds []uuid.UUID) (map[uuid.UUID]ledgerQuery.Period, error) {
	return s.ledgerInterface.ReadPeriodsByIds(ctx, periodIds)
}
