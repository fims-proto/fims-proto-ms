package domain

import (
	"time"
)

func (v *Voucher) UpdateLineItems(items []*LineItem) error {
	totalVal, err := sumItems(items)
	if err != nil {
		return err
	}
	if v.IsAudited() {
		return newDomainErr(errUpdateAudited)
	}
	if v.IsReviewed() {
		return newDomainErr(errUpdateReviewed)
	}
	// no need to check if voucher posted

	v.credit = totalVal
	v.debit = totalVal
	v.lineItems = items
	return nil
}

func (v *Voucher) UpdateTransactionTime(transactionTime time.Time) error {
	if transactionTime.IsZero() {
		return newDomainErr(errUpdateZeroTransactionTime)
	}
	if v.IsAudited() {
		return newDomainErr(errUpdateAudited)
	}
	if v.IsReviewed() {
		return newDomainErr(errUpdateReviewed)
	}
	// no need to check if voucher posted

	v.transactionTime = transactionTime
	return nil
}
