package journal_entry

import (
	"github.com/google/uuid"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (d *JournalEntry) Audit(auditor uuid.UUID) error {
	if d.isAudited {
		return commonErrors.NewSlugError("journalEntry-audit-repeatAudit", "entry is audited")
	}

	if auditor == uuid.Nil {
		return commonErrors.NewSlugError("journalEntry-audit-emptyAuditor", "empty auditor")
	}

	if auditor == d.creator {
		return commonErrors.NewSlugError("journalEntry-audit-auditorSameAsCreator", "auditor is same as creator")
	}

	if d.reviewer != uuid.Nil && auditor == d.reviewer {
		return commonErrors.NewSlugError("journalEntry-audit-auditorSameAsReviewer", "auditor is same as reviewer")
	}

	d.isAudited = true
	d.auditor = auditor
	return nil
}

func (d *JournalEntry) CancelAudit(auditor uuid.UUID) error {
	if !d.isAudited {
		return commonErrors.NewSlugError("journalEntry-cancelAudit-notAudited", "entry not audited")
	}

	if d.auditor != auditor {
		return commonErrors.NewSlugError("journalEntry-cancelAudit-differentAuditor", "only auditor can cancel audit")
	}

	if d.isPosted {
		return commonErrors.NewSlugError("journalEntry-cancelAudit-posted", "entry is posted")
	}

	d.isAudited = false
	d.auditor = uuid.Nil
	return nil
}
