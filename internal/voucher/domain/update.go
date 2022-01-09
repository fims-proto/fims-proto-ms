package domain

import (
	"time"

	"github.com/pkg/errors"
)

func (v *Voucher) UpdateLineItems(items []*LineItem) error {
	totalVal, err := sumItems(items)
	if err != nil {
		return err
	}
	if v.IsAudited() {
		return ErrVoucherAlreadyAudited
	}
	if v.IsReviewed() {
		return ErrVoucherAlreadyReviewed
	}
	if v.IsPosted() {
		return errors.New("voucher already posted")
	}

	v.credit = totalVal
	v.debit = totalVal
	v.lineItems = items
	return nil
}

func (v *Voucher) UpdateTransactionTime(transactionTime time.Time) error {
	if transactionTime.IsZero() {
		return errors.New("zero transaction time")
	}
	if v.IsAudited() {
		return ErrVoucherAlreadyAudited
	}
	if v.IsReviewed() {
		return ErrVoucherAlreadyReviewed
	}
	if v.IsPosted() {
		return errors.New("voucher already posted")
	}

	v.transactionTime = transactionTime
	return nil
}
