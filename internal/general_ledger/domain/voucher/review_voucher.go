package voucher

import (
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (d *Voucher) Review(reviewer uuid.UUID) error {
	if d.isReviewed {
		return errors.NewSlugError("voucher-review-repeatReview")
	}

	if reviewer == uuid.Nil {
		return errors.NewSlugError("voucher-review-emptyReviewer")
	}

	if reviewer == d.creator {
		return errors.NewSlugError("voucher-review-reviewerSameAsCreator")
	}

	if d.auditor != uuid.Nil && reviewer == d.auditor {
		return errors.NewSlugError("voucher-review-reviewerSameAsAuditor")
	}

	d.isReviewed = true
	d.reviewer = reviewer
	return nil
}

func (d *Voucher) CancelReview(reviewer uuid.UUID) error {
	if !d.isReviewed {
		return errors.NewSlugError("voucher-cancelReview-notReviewed")
	}

	if d.reviewer != reviewer {
		return errors.NewSlugError("voucher-cancelReview-differentReviewer")
	}

	if d.isPosted {
		return errors.NewSlugError("voucher-cancelReview-posted")
	}

	d.isReviewed = false
	d.reviewer = uuid.Nil
	return nil
}
