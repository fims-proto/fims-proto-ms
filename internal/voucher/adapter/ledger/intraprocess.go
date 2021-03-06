package ledger

import (
	"context"
	ledgerport "github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
)

type IntraprocessService struct {
	ledgerInterface ledgerport.LedgerInterface
}

func NewIntraprocessService(ledgerInterface ledgerport.LedgerInterface) IntraprocessService {
	return IntraprocessService{ledgerInterface: ledgerInterface}
}

func (s IntraprocessService) PostVoucher(ctx context.Context, voucher query.Voucher) error
