package voucher

import (
	"time"

	"github.com/pkg/errors"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/line_item"

	"github.com/google/uuid"
)

func (d *Voucher) checkUpdatePossible(user uuid.UUID) error {
	if d.isAudited {
		return commonErrors.NewSlugError("voucher-update-audited", "voucher is audited")
	}

	if d.isReviewed {
		return commonErrors.NewSlugError("voucher-update-reviewed", "voucher is reviewed")
	}

	if user != d.creator {
		return commonErrors.NewSlugError("voucher-update-notCreator", "only creator can update voucher")
	}

	return nil
}

func (d *Voucher) UpdateLineItems(lineItems []line_item.LineItem, user uuid.UUID) error {
	if err := d.checkUpdatePossible(user); err != nil {
		return errors.Wrap(err, "update not allowed")
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

func (d *Voucher) UpdateTransactionTime(transactionTime time.Time, periodId uuid.UUID, user uuid.UUID) error {
	if err := d.checkUpdatePossible(user); err != nil {
		return errors.Wrap(err, "update not allowed")
	}

	if periodId == uuid.Nil {
		return commonErrors.NewSlugError("voucher-emptyPeriodId", "empty period id")
	}

	if transactionTime.IsZero() {
		return commonErrors.NewSlugError("voucher-zeroTransactionTime", "empty period id")
	}

	if transactionTime.After(time.Now()) {
		return commonErrors.NewSlugError("voucher-futureTransactionTime", "empty period id")
	}

	d.transactionTime = transactionTime
	return nil
}

func (d *Voucher) UpdateHeaderText(headerText string, user uuid.UUID) error {
	if err := d.checkUpdatePossible(user); err != nil {
		return errors.Wrap(err, "update not allowed")
	}

	if headerText == "" {
		return errors.New("header text cannot be empty")
	}

	d.headerText = headerText
	return nil
}
