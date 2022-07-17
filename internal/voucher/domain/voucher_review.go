package domain

import "github.com/google/uuid"

func (v *Voucher) Review(reviewer uuid.UUID) error {
	if v.IsReviewed() {
		return newDomainErr(errReviewRepeatReview)
	}
	if reviewer == uuid.Nil {
		return newDomainErr(errReviewEmptyReviewer)
	}
	if reviewer == v.creator {
		return newDomainErr(errReviewReviewerSameAsCreator)
	}
	if v.auditor != uuid.Nil && reviewer == v.auditor {
		return newDomainErr(errReviewReviewerSameAsAuditor)
	}
	v.isReviewed = true
	v.reviewer = reviewer
	return nil
}

func (v *Voucher) CancelReview(reviewer uuid.UUID) error {
	if !v.IsReviewed() {
		return newDomainErr(errCancelReviewNotReviewed)
	}
	if v.Reviewer() != reviewer {
		return newDomainErr(errCancelReviewDifferentReviewer)
	}
	if v.IsPosted() {
		return newDomainErr(errCancelReviewPosted)
	}
	v.isReviewed = false
	v.reviewer = uuid.Nil
	return nil
}
