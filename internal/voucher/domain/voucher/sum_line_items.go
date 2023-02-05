package voucher

import (
	"github.com/shopspring/decimal"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/line_item"
)

func sumLineItems(lineItems []line_item.LineItem) (decimal.Decimal, error) {
	if len(lineItems) == 0 {
		return decimal.Decimal{}, commonErrors.NewSlugError("voucher-emptyLineItems", "empty lineItems")
	}

	var debitInTotal decimal.Decimal
	var creditInTotal decimal.Decimal
	for _, item := range lineItems {
		debitInTotal = debitInTotal.Add(item.Debit())
		creditInTotal = creditInTotal.Add(item.Credit())
	}

	if !debitInTotal.Equal(creditInTotal) {
		return decimal.Decimal{}, commonErrors.NewSlugError("voucher-notBalanced", "voucher is not balanced")
	}
	return debitInTotal, nil
}
