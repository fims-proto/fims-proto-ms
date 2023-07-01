package voucher

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (d *Voucher) checkUpdatePossible(user uuid.UUID) error {
	if d.isAudited {
		return commonErrors.NewSlugError("voucher-update-audited")
	}

	if d.isReviewed {
		return commonErrors.NewSlugError("voucher-update-reviewed")
	}

	if user != d.creator {
		return commonErrors.NewSlugError("voucher-update-notCreator")
	}

	return nil
}

func (d *Voucher) UpdateLineItems(lineItems []*LineItem, user uuid.UUID) error {
	if err := d.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	totalVal, err := sumLineItems(lineItems)
	if err != nil {
		return err
	}

	d.credit = totalVal
	d.debit = totalVal
	d.lineItems = lineItems
	return nil
}

func (d *Voucher) UpdateTransactionTime(transactionTime time.Time, user uuid.UUID) error {
	if err := d.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if transactionTime.IsZero() {
		return commonErrors.NewSlugError("voucher-zeroTransactionTime")
	}

	d.transactionTime = transactionTime
	return nil
}

func (d *Voucher) UpdatePeriodAndDocumentNumber(periodId uuid.UUID, documentNumber string, user uuid.UUID) error {
	if err := d.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if periodId == uuid.Nil {
		return commonErrors.NewSlugError("voucher-emptyPeriodId")
	}

	if documentNumber == "" {
		return commonErrors.NewSlugError("voucher-emptyNumber")
	}

	d.periodId = periodId
	d.documentNumber = documentNumber
	return nil
}

func (d *Voucher) UpdateHeaderText(headerText string, user uuid.UUID) error {
	if err := d.checkUpdatePossible(user); err != nil {
		return fmt.Errorf("update not allowed: %w", err)
	}

	if headerText == "" {
		return commonErrors.NewSlugError("voucher-emptyHeaderText")
	}

	d.headerText = headerText
	return nil
}
