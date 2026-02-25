package voucher

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/shopspring/decimal"
)

func sumLineItems(lineItems []*LineItem) (decimal.Decimal, error) {
	if len(lineItems) == 0 {
		return decimal.Decimal{}, errors.NewSlugError("voucher-emptyLineItems")
	}

	var sumAmount decimal.Decimal
	for _, item := range lineItems {
		if item == nil {
			return decimal.Decimal{}, errors.NewSlugError("voucher-nilLineItem")
		}

		sumAmount = sumAmount.Add(item.Amount())
	}

	// Trial balance: sum of all signed amounts must be zero
	if !sumAmount.IsZero() {
		return decimal.Decimal{}, errors.NewSlugError("voucher-notBalanced")
	}

	// Return the transaction amount (sum of all positive amounts/debits only)
	var transactionAmount decimal.Decimal
	for _, item := range lineItems {
		if item.Amount().IsPositive() {
			transactionAmount = transactionAmount.Add(item.Amount())
		}
	}

	return transactionAmount, nil
}
