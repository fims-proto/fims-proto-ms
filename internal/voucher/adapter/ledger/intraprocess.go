package ledger

import (
	"context"
	ledgerport "github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
)

type IntraprocessAdapter struct {
	ledgerInterface ledgerport.LedgerInterface
}

func NewIntraprocessAdapter(ledgerInterface ledgerport.LedgerInterface) IntraprocessAdapter {
	return IntraprocessAdapter{ledgerInterface: ledgerInterface}
}

func (s IntraprocessAdapter) PostVoucher(ctx context.Context, voucher query.Voucher) error {
	return s.ledgerInterface.PostVoucher(ctx, mapFromVoucherQuery(voucher))
}
