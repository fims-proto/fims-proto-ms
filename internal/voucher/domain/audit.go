package domain

import "github.com/google/uuid"

func (v *Voucher) Audit(auditor uuid.UUID) error {
	if v.IsAudited() {
		return newDomainErr(errAuditRepeatAudit)
	}
	if auditor == uuid.Nil {
		return newDomainErr(errAuditEmptyAuditor)
	}
	if auditor == v.creator {
		return newDomainErr(errAuditAuditorSameAsCreator)
	}
	if v.reviewer != uuid.Nil && auditor == v.reviewer {
		return newDomainErr(errAuditAuditorSameAsReviewer)
	}
	v.isAudited = true
	v.auditor = auditor
	return nil
}

func (v *Voucher) CancelAudit(auditor uuid.UUID) error {
	if !v.IsAudited() {
		return newDomainErr(errCancelAuditNotAudited)
	}
	if v.Auditor() != auditor {
		return newDomainErr(errCancelAuditDifferentAuditor)
	}
	if v.IsPosted() {
		return newDomainErr(errCancelAuditPosted)
	}
	v.isAudited = false
	v.auditor = uuid.Nil
	return nil
}
