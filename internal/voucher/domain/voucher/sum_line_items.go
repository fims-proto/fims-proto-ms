package voucher

import (
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/line_item"
)

func sumLineItems(lineItems []line_item.LineItem) (decimal.Decimal, error) {
	if len(lineItems) == 0 {
		return decimal.Decimal{}, errors.NewSlugError("voucher-emptyLineItems")
	}

	var debitInTotal decimal.Decimal
	var creditInTotal decimal.Decimal
	for _, item := range lineItems {
		debitInTotal = debitInTotal.Add(item.Debit())
		creditInTotal = creditInTotal.Add(item.Credit())
	}

	if !debitInTotal.Equal(creditInTotal) {
		return decimal.Decimal{}, errors.NewSlugError("voucher-notBalanced")
	}
	return debitInTotal, nil
}
