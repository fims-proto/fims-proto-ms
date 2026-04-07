package journal

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (j *Journal) Audit(auditor string) error {
	if j.isAudited {
		return errors.NewSlugError("journal-audit-repeatAudit")
	}

	if isEmptyUser(auditor) {
		return errors.NewSlugError("journal-audit-emptyAuditor")
	}

	if !IsSystemUser(auditor) {
		if auditor == j.creator {
			return errors.NewSlugError("journal-audit-auditorSameAsCreator")
		}

		if !isEmptyUser(j.reviewer) && auditor == j.reviewer {
			return errors.NewSlugError("journal-audit-auditorSameAsReviewer")
		}
	}

	j.isAudited = true
	j.auditor = auditor
	return nil
}

func (j *Journal) CancelAudit(auditor string) error {
	if !j.isAudited {
		return errors.NewSlugError("journal-cancelAudit-notAudited")
	}

	if j.auditor != auditor {
		return errors.NewSlugError("journal-cancelAudit-differentAuditor")
	}

	if j.isPosted {
		return errors.NewSlugError("journal-cancelAudit-posted")
	}

	j.isAudited = false
	j.auditor = emptyUser
	return nil
}
