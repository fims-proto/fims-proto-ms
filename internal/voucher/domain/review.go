package domain

import "github.com/pkg/errors"

var (
	ErrEmptyReviewer           = errors.New("reviewer empty")
	ErrVoucherAlreadyReviewed  = errors.New("voucher already reviewed")
	ErrVoucherNotReviewed      = errors.New("voucher not reviewed")
	ErrDifferentReviewerCancel = errors.New("cancel review with different reviewer")
)

func (v *Voucher) Review(reviewer string) error {
	if v.IsReviewed() {
		return ErrVoucherAlreadyReviewed
	}
	if reviewer == "" {
		return ErrEmptyReviewer
	}
	if reviewer == v.creator {
		return errors.New("reviewer cannot be same as creator")
	}
	if v.auditor != "" && reviewer == v.auditor {
		return errors.New("reviewer cannot be same as auditor")
	}
	v.isReviewed = true
	v.reviewer = reviewer
	return nil
}

func (v *Voucher) CancelReview(reviewer string) error {
	if !v.IsReviewed() {
		return ErrVoucherNotReviewed
	}
	if v.Reviewer() != reviewer {
		return ErrDifferentReviewerCancel
	}
	v.isReviewed = false
	v.reviewer = ""
	return nil
}
