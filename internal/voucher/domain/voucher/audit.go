package voucher

import "github.com/pkg/errors"

var (
	ErrVoucherAlreadyAudited  = errors.New("voucher already audited")
	ErrVoucherNotAudited      = errors.New("voucher not audited")
	ErrDifferentAuditorCancel = errors.New("cancel audit with different auditor")
)

func (v *Voucher) Audit(supervisorUUID string) error {
	if v.auditor.isAudited {
		return ErrVoucherAlreadyAudited
	}
	v.auditor.isAudited = true
	v.auditor.uuid = supervisorUUID
	return nil
}

func (v *Voucher) CancelAudit(supervisorUUID string) error {
	if !v.auditor.isAudited {
		return ErrVoucherNotAudited
	}
	if v.auditor.uuid != supervisorUUID {
		return ErrDifferentAuditorCancel
	}
	v.auditor.isAudited = false
	v.auditor.uuid = ""
	return nil
}