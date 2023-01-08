package journal_entry

import (
	"github.com/google/uuid"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (d *JournalEntry) Review(reviewer uuid.UUID) error {
	if d.isReviewed {
		return commonErrors.NewSlugError("journalEntry-review-repeatReview", "entry is reviewed")
	}

	if reviewer == uuid.Nil {
		return commonErrors.NewSlugError("journalEntry-review-emptyReviewer", "empty reviewer")
	}

	if reviewer == d.creator {
		return commonErrors.NewSlugError("journalEntry-review-reviewerSameAsCreator", "reviewer is same as creator")
	}

	if d.auditor != uuid.Nil && reviewer == d.auditor {
		return commonErrors.NewSlugError("journalEntry-review-reviewerSameAsAuditor", "reviewer is same as auditor")
	}

	d.isReviewed = true
	d.reviewer = reviewer
	return nil
}

func (d *JournalEntry) CancelReview(reviewer uuid.UUID) error {
	if !d.isReviewed {
		return commonErrors.NewSlugError("journalEntry-cancelReview-notReviewed", "entry is not reviewed")
	}

	if d.reviewer != reviewer {
		return commonErrors.NewSlugError("journalEntry-cancelReview-differentReviewer", "only reviewer can cancel review")
	}

	if d.isPosted {
		return commonErrors.NewSlugError("journalEntry-cancelReview-posted", "entry is posted")
	}

	d.isReviewed = false
	d.reviewer = uuid.Nil
	return nil
}
