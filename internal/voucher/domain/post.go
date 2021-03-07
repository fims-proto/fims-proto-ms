package domain

import "github.com/pkg/errors"

func (v *Voucher) Post() error {
	if v.IsPosted() {
		return errors.Errorf("voucher %s already posted", v.Number())
	}
	v.isPosted = true
	return nil
}
