package journal

import (
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

func (j *Journal) checkUpdatePossible(user string) error {
	if j.isAudited {
		return commonErrors.NewInvalidInputError(commonErrors.SlugJournalUpdateAudited)
	}

	if j.isReviewed {
		return commonErrors.NewInvalidInputError(commonErrors.SlugJournalUpdateReviewed)
	}

	if !IsSystemUser(user) && user != j.creator {
		return commonErrors.NewInvalidInputError(commonErrors.SlugJournalUpdateNotCreator)
	}

	return nil
}

func (j *Journal) UpdateJournalLines(journalLines []*JournalLine, user string) error {
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

func (j *Journal) UpdateTransactionDate(transactionDate transaction_date.TransactionDate, user string) error {
	if err := j.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if transactionDate.IsZero() {
		return commonErrors.NewInvalidInputError(commonErrors.SlugJournalZeroTransactionDate)
	}

	j.transactionDate = transactionDate
	return nil
}

func (j *Journal) UpdatePeriodAndDocumentNumber(periodId uuid.UUID, documentNumber string, user string) error {
	if err := j.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if periodId == uuid.Nil {
		return commonErrors.NewInternalError(commonErrors.SlugJournalEmptyPeriodId)
	}

	if documentNumber == "" {
		return commonErrors.NewInvalidInputError(commonErrors.SlugJournalEmptyNumber)
	}

	j.periodId = periodId
	j.documentNumber = documentNumber
	return nil
}

func (j *Journal) UpdateHeaderText(headerText string, user string) error {
	if err := j.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if headerText == "" {
		return commonErrors.NewInvalidInputError(commonErrors.SlugJournalEmptyHeaderText)
	}

	j.headerText = headerText
	return nil
}
