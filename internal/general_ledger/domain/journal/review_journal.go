package journal

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (j *Journal) Review(reviewer string) error {
	if j.isReviewed {
		return errors.NewSlugError("journal-review-repeatReview")
	}

	if isEmptyUser(reviewer) {
		return errors.NewSlugError("journal-review-emptyReviewer")
	}

	if !IsSystemUser(reviewer) {
		if reviewer == j.creator {
			return errors.NewSlugError("journal-review-reviewerSameAsCreator")
		}

		if !isEmptyUser(j.auditor) && reviewer == j.auditor {
			return errors.NewSlugError("journal-review-reviewerSameAsAuditor")
		}
	}

	j.isReviewed = true
	j.reviewer = reviewer
	return nil
}

func (j *Journal) CancelReview(reviewer string) error {
	if !j.isReviewed {
		return errors.NewSlugError("journal-cancelReview-notReviewed")
	}

	if j.reviewer != reviewer {
		return errors.NewSlugError("journal-cancelReview-differentReviewer")
	}

	if j.isPosted {
		return errors.NewSlugError("journal-cancelReview-posted")
	}

	j.isReviewed = false
	j.reviewer = emptyUser
	return nil
}
