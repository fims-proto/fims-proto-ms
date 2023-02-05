package voucher

import (
	"github.com/google/uuid"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (d *Voucher) Review(reviewer uuid.UUID) error {
	if d.isReviewed {
		return commonErrors.NewSlugError("voucher-review-repeatReview", "voucher is reviewed")
	}

	if reviewer == uuid.Nil {
		return commonErrors.NewSlugError("voucher-review-emptyReviewer", "empty reviewer")
	}

	if reviewer == d.creator {
		return commonErrors.NewSlugError("voucher-review-reviewerSameAsCreator", "reviewer is same as creator")
	}

	if d.auditor != uuid.Nil && reviewer == d.auditor {
		return commonErrors.NewSlugError("voucher-review-reviewerSameAsAuditor", "reviewer is same as auditor")
	}

	d.isReviewed = true
	d.reviewer = reviewer
	return nil
}

func (d *Voucher) CancelReview(reviewer uuid.UUID) error {
	if !d.isReviewed {
		return commonErrors.NewSlugError("voucher-cancelReview-notReviewed", "voucher is not reviewed")
	}

	if d.reviewer != reviewer {
		return commonErrors.NewSlugError("voucher-cancelReview-differentReviewer", "only reviewer can cancel review")
	}

	if d.isPosted {
		return commonErrors.NewSlugError("voucher-cancelReview-posted", "voucher is posted")
	}

	d.isReviewed = false
	d.reviewer = uuid.Nil
	return nil
}
