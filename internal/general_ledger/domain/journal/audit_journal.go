package journal

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

func (j *Journal) Audit(auditor uuid.UUID) error {
	if j.isAudited {
		return errors.NewSlugError("journal-audit-repeatAudit")
	}

	if auditor == uuid.Nil {
		return errors.NewSlugError("journal-audit-emptyAuditor")
	}

	if auditor == j.creator {
		return errors.NewSlugError("journal-audit-auditorSameAsCreator")
	}

	if j.reviewer != uuid.Nil && auditor == j.reviewer {
		return errors.NewSlugError("journal-audit-auditorSameAsReviewer")
	}

	j.isAudited = true
	j.auditor = auditor
	return nil
}

func (j *Journal) CancelAudit(auditor uuid.UUID) error {
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
	j.auditor = uuid.Nil
	return nil
}
