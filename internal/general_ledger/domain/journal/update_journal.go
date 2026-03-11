package journal

import (
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

func (j *Journal) checkUpdatePossible(user uuid.UUID) error {
	if j.isAudited {
		return commonErrors.NewSlugError("journal-update-audited")
	}

	if j.isReviewed {
		return commonErrors.NewSlugError("journal-update-reviewed")
	}

	if user != j.creator {
		return commonErrors.NewSlugError("journal-update-notCreator")
	}

	return nil
}

func (j *Journal) UpdateJournalLines(journalLines []*JournalLine, user uuid.UUID) error {
	if err := j.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	totalVal, err := sumJournalLines(journalLines)
	if err != nil {
		return err
	}

	j.amount = totalVal
	j.journalLines = journalLines
	return nil
}

func (j *Journal) UpdateTransactionDate(transactionDate transaction_date.TransactionDate, user uuid.UUID) error {
	if err := j.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if transactionDate.IsZero() {
		return commonErrors.NewSlugError("journal-zeroTransactionDate")
	}

	j.transactionDate = transactionDate
	return nil
}

func (j *Journal) UpdatePeriodAndDocumentNumber(periodId uuid.UUID, documentNumber string, user uuid.UUID) error {
	if err := j.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if periodId == uuid.Nil {
		return commonErrors.NewSlugError("journal-emptyPeriodId")
	}

	if documentNumber == "" {
		return commonErrors.NewSlugError("journal-emptyNumber")
	}

	j.periodId = periodId
	j.documentNumber = documentNumber
	return nil
}

func (j *Journal) UpdateHeaderText(headerText string, user uuid.UUID) error {
	if err := j.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if headerText == "" {
		return commonErrors.NewSlugError("journal-emptyHeaderText")
	}

	j.headerText = headerText
	return nil
}
