package ledger

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
	ledgerport "github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"
)

type IntraprocessAdapter struct {
	ledgerInterface ledgerport.LedgerInterface
}

func NewIntraprocessAdapter(ledgerInterface ledgerport.LedgerInterface) IntraprocessAdapter {
	return IntraprocessAdapter{ledgerInterface: ledgerInterface}
}

func (s IntraprocessAdapter) LoadLedgers(ctx context.Context, accounts []domain.Account) error {
	var reqs []ledgerport.LoadLedgersRequest
	for _, account := range accounts {
		reqs = append(reqs, mapFromDomainAccount(account))
	}
	return s.ledgerInterface.LoadLedgers(ctx, reqs)
}
