package domain

func (v *Voucher) Audit(auditor string) error {
	if v.IsAudited() {
		return newDomainErr(errAuditRepeatAudit)
	}
	if auditor == "" {
		return newDomainErr(errAuditEmptyAuditor)
	}
	if auditor == v.creator {
		return newDomainErr(errAuditAuditorSameAsCreator)
	}
	if v.reviewer != "" && auditor == v.reviewer {
		return newDomainErr(errAuditAuditorSameAsReviewer)
	}
	v.isAudited = true
	v.auditor = auditor
	return nil
}

func (v *Voucher) CancelAudit(auditor string) error {
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
	v.auditor = ""
	return nil
}
