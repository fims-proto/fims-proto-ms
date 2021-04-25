package ledger

import (
	"github/fims-proto/fims-proto-ms/internal/account/domain"
	ledgerport "github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"
)

func mapFromDomainAccount(a domain.Account) ledgerport.LoadLedgersRequest {
	return ledgerport.LoadLedgersRequest{
		Number:         a.Number(),
		Title:          a.Title(),
		SuperiorNumber: a.SuperiorNumber(),
		AccountType:    a.Type(),
	}
}
