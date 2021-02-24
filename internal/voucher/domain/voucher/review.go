package voucher

import "github.com/pkg/errors"

var (
	ErrEmptyReviewer           = errors.New("reviewer empty")
	ErrVoucherAlreadyReviewed  = errors.New("voucher already reiviewed")
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
