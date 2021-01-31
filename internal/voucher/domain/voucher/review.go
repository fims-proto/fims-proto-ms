package voucher

import "github.com/pkg/errors"

var (
	ErrEmptyReviewer           = errors.New("reviewer uuid empty")
	ErrVoucherAlreadyReviewed  = errors.New("voucher already audited")
	ErrVoucherNotReviewed      = errors.New("voucher not reviewed")
	ErrDifferentReviewerCancel = errors.New("cancel review with different reviewer")
)

func (v *Voucher) Review(reviewerUUID string) error {
	if v.reviewer.isReviewed {
		return ErrVoucherAlreadyReviewed
	}
	if reviewerUUID == "" {
		return ErrEmptyReviewer
	}
	v.reviewer.isReviewed = true
	v.reviewer.uuid = reviewerUUID
	return nil
}

func (v *Voucher) CancelReview(reviewerUUID string) error {
	if !v.reviewer.isReviewed {
		return ErrVoucherNotReviewed
	}
	if v.reviewer.uuid != reviewerUUID {
		return ErrDifferentReviewerCancel
	}
	v.reviewer.isReviewed = false
	v.reviewer.uuid = ""
	return nil
}
