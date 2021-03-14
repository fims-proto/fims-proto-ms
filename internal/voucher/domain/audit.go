package domain

import "github.com/pkg/errors"

var (
	ErrEmptyAuditor           = errors.New("auditor empty")
	ErrVoucherAlreadyAudited  = errors.New("voucher already audited")
	ErrVoucherNotAudited      = errors.New("voucher not audited")
	ErrDifferentAuditorCancel = errors.New("cancel audit with different auditor")
)

func (v *Voucher) Audit(auditor string) error {
	if v.IsAudited() {
		return ErrVoucherAlreadyAudited
	}
	if auditor == "" {
		return ErrEmptyAuditor
	}
	v.isAudited = true
	v.auditor = auditor
	return nil
}

func (v *Voucher) CancelAudit(auditor string) error {
	if !v.IsAudited() {
		return ErrVoucherNotAudited
	}
	if v.Auditor() != auditor {
		return ErrDifferentAuditorCancel
	}
	v.isAudited = false
	v.auditor = ""
	return nil
}
