package voucher

import "github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"

func (v *Voucher) Update(items []lineitem.LineItem) error {
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

	v.credit = totalVal
	v.debit = totalVal
	v.lineItems = items
	return nil
}
