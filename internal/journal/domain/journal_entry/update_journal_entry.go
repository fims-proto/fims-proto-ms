package journal_entry

import (
	"time"

	"github.com/pkg/errors"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/line_item"

	"github.com/google/uuid"
)

func (d *JournalEntry) checkUpdatePossible(user uuid.UUID) error {
	if d.isAudited {
		return commonErrors.NewSlugError("journalEntry-update-audited", "entry is audited")
	}

	if d.isReviewed {
		return commonErrors.NewSlugError("journalEntry-update-reviewed", "entry is reviewed")
	}

	if user != d.creator {
		return commonErrors.NewSlugError("journalEntry-update-notCreator", "only creator can update entry")
	}

	return nil
}

func (d *JournalEntry) UpdateLineItems(lineItems []line_item.LineItem, user uuid.UUID) error {
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

func (d *JournalEntry) UpdateTransactionTime(transactionTime time.Time, periodId uuid.UUID, user uuid.UUID) error {
	if err := d.checkUpdatePossible(user); err != nil {
		return errors.Wrap(err, "update not allowed")
	}

	if periodId == uuid.Nil {
		return commonErrors.NewSlugError("journalEntry-emptyPeriodId", "empty period id")
	}

	if transactionTime.IsZero() {
		return commonErrors.NewSlugError("journalEntry-zeroTransactionTime", "empty period id")
	}

	if transactionTime.After(time.Now()) {
		return commonErrors.NewSlugError("journalEntry-futureTransactionTime", "empty period id")
	}

	d.transactionTime = transactionTime
	return nil
}

func (d *JournalEntry) UpdateHeaderText(headerText string, user uuid.UUID) error {
	if err := d.checkUpdatePossible(user); err != nil {
		return errors.Wrap(err, "update not allowed")
	}

	if headerText == "" {
		return errors.New("header text cannot be empty")
	}

	d.headerText = headerText
	return nil
}
