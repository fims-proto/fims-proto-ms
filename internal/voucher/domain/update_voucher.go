package domain

import (
	"time"

	"github.com/google/uuid"
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

func (v *Voucher) UpdateTransactionTime(transactionTime time.Time, periodId uuid.UUID) error {
	if periodId == uuid.Nil {
		return newDomainErr(errVoucherEmptyPeriodId)
	}
	if transactionTime.IsZero() {
		return newDomainErr(errUpdateZeroTransactionTime)
	}
	if transactionTime.After(time.Now()) {
		return newDomainErr(errVoucherFutureTransactionTime, transactionTime)
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
