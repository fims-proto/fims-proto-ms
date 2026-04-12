package journal

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (j *Journal) Review(reviewer string) error {
	if j.isReviewed {
		return errors.NewInvalidInputError(errors.SlugJournalReviewRepeat)
	}

	if isEmptyUser(reviewer) {
		return errors.NewInternalError(errors.SlugJournalReviewEmptyReviewer)
	}

	if !IsSystemUser(reviewer) {
		if reviewer == j.creator {
			return errors.NewInvalidInputError(errors.SlugJournalReviewSameAsCreator)
		}

		if !isEmptyUser(j.auditor) && reviewer == j.auditor {
			return errors.NewInvalidInputError(errors.SlugJournalReviewSameAsAuditor)
		}
	}

	j.isReviewed = true
	j.reviewer = reviewer
	return nil
}

func (j *Journal) CancelReview(reviewer string) error {
	if !j.isReviewed {
		return errors.NewInvalidInputError(errors.SlugJournalCancelReviewNotReviewed)
	}

	if j.reviewer != reviewer {
		return errors.NewInvalidInputError(errors.SlugJournalCancelReviewDiffReviewer)
	}

	if j.isPosted {
		return errors.NewInvalidInputError(errors.SlugJournalCancelReviewPosted)
	}

	j.isReviewed = false
	j.reviewer = emptyUser
	return nil
}
