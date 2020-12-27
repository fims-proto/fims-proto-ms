package voucher

import "github.com/pkg/errors"

var (
	ErrVoucherAlreadyAudited  = errors.New("voucher already audited")
	ErrVoucherNotAudited      = errors.New("voucher not audited")
	ErrDifferentAuditorCancel = errors.New("cancel audit with different auditor")
)

func (v *Voucher) Audit(auditorUUID string) error {
	if v.auditor.isAudited {
		return ErrVoucherAlreadyAudited
	}
	v.auditor.isAudited = true
	v.auditor.uuid = auditorUUID
	return nil
}

func (v *Voucher) CancelAudit(auditorUUID string) error {
	if !v.auditor.isAudited {
		return ErrVoucherNotAudited
	}
	if v.auditor.uuid != auditorUUID {
		return ErrDifferentAuditorCancel
	}
	v.auditor.isAudited = false
	v.auditor.uuid = ""
	return nil
}