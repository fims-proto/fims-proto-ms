package domain

func (v *Voucher) Review(reviewer string) error {
	if v.IsReviewed() {
		return newDomainErr(errReviewRepeatReview)
	}
	if reviewer == "" {
		return newDomainErr(errReviewEmptyReviewer)
	}
	if reviewer == v.creator {
		return newDomainErr(errReviewReviewerSameAsCreator)
	}
	if v.auditor != "" && reviewer == v.auditor {
		return newDomainErr(errReviewReviewerSameAsAuditor)
	}
	v.isReviewed = true
	v.reviewer = reviewer
	return nil
}

func (v *Voucher) CancelReview(reviewer string) error {
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
	v.reviewer = ""
	return nil
}
