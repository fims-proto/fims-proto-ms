package ledger

import (
	ledgerport "github/fims-proto/fims-proto-ms/internal/ledger/port/private/intraprocess"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"

	"github.com/shopspring/decimal"
)

func mapFromLineItemQuery(q query.LineItem) ledgerport.VoucherLineItemRequest {
	return ledgerport.VoucherLineItemRequest{
		AccountNumber: q.AccountNumber,
		Debit:         decimal.RequireFromString(q.Debit),
		Credit:        decimal.RequireFromString(q.Credit),
	}
}

func mapFromVoucherQuery(q query.Voucher) ledgerport.UpdateLedgerBalanceRequest {
	itemReq := []ledgerport.VoucherLineItemRequest{}
	for _, item := range q.LineItems {
		itemReq = append(itemReq, mapFromLineItemQuery(item))
	}
	return ledgerport.UpdateLedgerBalanceRequest{
		Sob:         q.Sob,
		VoucherUUID: q.UUID,
		LineItems:   itemReq,
	}
}
