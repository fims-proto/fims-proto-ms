package voucher

import (
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (v *Voucher) Review(reviewer uuid.UUID) error {
	if v.isReviewed {
		return errors.NewSlugError("voucher-review-repeatReview")
	}

	if reviewer == uuid.Nil {
		return errors.NewSlugError("voucher-review-emptyReviewer")
	}

	if reviewer == v.creator {
		return errors.NewSlugError("voucher-review-reviewerSameAsCreator")
	}

	if v.auditor != uuid.Nil && reviewer == v.auditor {
		return errors.NewSlugError("voucher-review-reviewerSameAsAuditor")
	}

	v.isReviewed = true
	v.reviewer = reviewer
	return nil
}

func (v *Voucher) CancelReview(reviewer uuid.UUID) error {
	if !v.isReviewed {
		return errors.NewSlugError("voucher-cancelReview-notReviewed")
	}

	if v.reviewer != reviewer {
		return errors.NewSlugError("voucher-cancelReview-differentReviewer")
	}

	if v.isPosted {
		return errors.NewSlugError("voucher-cancelReview-posted")
	}

	v.isReviewed = false
	v.reviewer = uuid.Nil
	return nil
}
