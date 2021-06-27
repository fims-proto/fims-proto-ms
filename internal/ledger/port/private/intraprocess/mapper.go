package intraprocess

import "github/fims-proto/fims-proto-ms/internal/ledger/app/command"

func (r UpdateLedgerBalanceRequest) mapToCommand() command.UpdateLedgerBalanceCmd {
	items := []command.LineItemCmd{}
	for _, i := range r.LineItems {
		items = append(items, i.mapToCommand())
	}
	return command.UpdateLedgerBalanceCmd{
		Sob:         r.Sob,
		VoucherUUID: r.VoucherUUID,
		LineItems:   items,
	}
}

func (r VoucherLineItemRequest) mapToCommand() command.LineItemCmd {
	return command.LineItemCmd{
		AccountNumber: r.AccountNumber,
		Debit:         r.Debit,
		Credit:        r.Credit,
	}
}

func (r LoadLedgersRequest) mapToCommand() command.LedgerDataloadCmd {
	return command.LedgerDataloadCmd{
		Number:         r.Number,
		Title:          r.Title,
		SuperiorNumber: r.SuperiorNumber,
		AccountType:    r.AccountType,
	}
}
