package journal

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (j *Journal) Audit(auditor string) error {
	if j.isAudited {
		return errors.NewInvalidInputError(errors.SlugJournalAuditRepeatAudit)
	}

	if isEmptyUser(auditor) {
		return errors.NewInternalError(errors.SlugJournalAuditEmptyAuditor)
	}

	if !IsSystemUser(auditor) {
		if auditor == j.creator {
			return errors.NewInvalidInputError(errors.SlugJournalAuditSameAsCreator)
		}

		if !isEmptyUser(j.reviewer) && auditor == j.reviewer {
			return errors.NewInvalidInputError(errors.SlugJournalAuditSameAsReviewer)
		}
	}

	j.isAudited = true
	j.auditor = auditor
	return nil
}

func (j *Journal) CancelAudit(auditor string) error {
	if !j.isAudited {
		return errors.NewInvalidInputError(errors.SlugJournalCancelAuditNotAudited)
	}

	if j.auditor != auditor {
		return errors.NewInvalidInputError(errors.SlugJournalCancelAuditDiffAuditor)
	}

	if j.isPosted {
		return errors.NewInvalidInputError(errors.SlugJournalCancelAuditPosted)
	}

	j.isAudited = false
	j.auditor = emptyUser
	return nil
}
