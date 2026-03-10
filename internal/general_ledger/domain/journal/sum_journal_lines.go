package journal

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/shopspring/decimal"
)

func sumJournalLines(journalLines []*JournalLine) (decimal.Decimal, error) {
	if len(journalLines) == 0 {
		return decimal.Decimal{}, errors.NewSlugError("journal-emptyJournalLines")
	}

	var sumAmount decimal.Decimal
	for _, item := range journalLines {
		if item == nil {
			return decimal.Decimal{}, errors.NewSlugError("journal-nilJournalLine")
		}

		sumAmount = sumAmount.Add(item.Amount())
	}

	// Trial balance: sum of all signed amounts must be zero
	if !sumAmount.IsZero() {
		return decimal.Decimal{}, errors.NewSlugError("journal-notBalanced")
	}

	// Return the transaction amount (sum of all positive amounts/debits only)
	var transactionAmount decimal.Decimal
	for _, item := range journalLines {
		if item.Amount().IsPositive() {
			transactionAmount = transactionAmount.Add(item.Amount())
		}
	}

	return transactionAmount, nil
}
