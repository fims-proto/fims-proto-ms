package domain

func (v *Voucher) Update(items []LineItem) error {
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
