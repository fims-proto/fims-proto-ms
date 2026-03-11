package journal

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

func (j *Journal) Review(reviewer uuid.UUID) error {
	if j.isReviewed {
		return errors.NewSlugError("journal-review-repeatReview")
	}

	if reviewer == uuid.Nil {
		return errors.NewSlugError("journal-review-emptyReviewer")
	}

	if reviewer == j.creator {
		return errors.NewSlugError("journal-review-reviewerSameAsCreator")
	}

	if j.auditor != uuid.Nil && reviewer == j.auditor {
		return errors.NewSlugError("journal-review-reviewerSameAsAuditor")
	}

	j.isReviewed = true
	j.reviewer = reviewer
	return nil
}

func (j *Journal) CancelReview(reviewer uuid.UUID) error {
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
	j.reviewer = uuid.Nil
	return nil
}
