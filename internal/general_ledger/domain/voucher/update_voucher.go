package voucher

import (
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

func (v *Voucher) checkUpdatePossible(user uuid.UUID) error {
	if v.isAudited {
		return commonErrors.NewSlugError("voucher-update-audited")
	}

	if v.isReviewed {
		return commonErrors.NewSlugError("voucher-update-reviewed")
	}

	if user != v.creator {
		return commonErrors.NewSlugError("voucher-update-notCreator")
	}

	return nil
}

func (v *Voucher) UpdateLineItems(lineItems []*LineItem, user uuid.UUID) error {
	if err := v.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	totalVal, err := sumLineItems(lineItems)
	if err != nil {
		return err
	}

	v.amount = totalVal
	v.lineItems = lineItems
	return nil
}

func (v *Voucher) UpdateTransactionDate(transactionDate transaction_date.TransactionDate, user uuid.UUID) error {
	if err := v.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if transactionDate.IsZero() {
		return commonErrors.NewSlugError("voucher-zeroTransactionDate")
	}

	v.transactionDate = transactionDate
	return nil
}

func (v *Voucher) UpdatePeriodAndDocumentNumber(periodId uuid.UUID, documentNumber string, user uuid.UUID) error {
	if err := v.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if periodId == uuid.Nil {
		return commonErrors.NewSlugError("voucher-emptyPeriodId")
	}

	if documentNumber == "" {
		return commonErrors.NewSlugError("voucher-emptyNumber")
	}

	v.periodId = periodId
	v.documentNumber = documentNumber
	return nil
}

func (v *Voucher) UpdateHeaderText(headerText string, user uuid.UUID) error {
	if err := v.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if headerText == "" {
		return commonErrors.NewSlugError("voucher-emptyHeaderText")
	}

	v.headerText = headerText
	return nil
}
