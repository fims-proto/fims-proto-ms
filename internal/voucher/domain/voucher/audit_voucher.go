package voucher

import (
	"github.com/google/uuid"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (d *Voucher) Audit(auditor uuid.UUID) error {
	if d.isAudited {
		return commonErrors.NewSlugError("voucher-audit-repeatAudit", "voucher is audited")
	}

	if auditor == uuid.Nil {
		return commonErrors.NewSlugError("voucher-audit-emptyAuditor", "empty auditor")
	}

	if auditor == d.creator {
		return commonErrors.NewSlugError("voucher-audit-auditorSameAsCreator", "auditor is same as creator")
	}

	if d.reviewer != uuid.Nil && auditor == d.reviewer {
		return commonErrors.NewSlugError("voucher-audit-auditorSameAsReviewer", "auditor is same as reviewer")
	}

	d.isAudited = true
	d.auditor = auditor
	return nil
}

func (d *Voucher) CancelAudit(auditor uuid.UUID) error {
	if !d.isAudited {
		return commonErrors.NewSlugError("voucher-cancelAudit-notAudited", "voucher not audited")
	}

	if d.auditor != auditor {
		return commonErrors.NewSlugError("voucher-cancelAudit-differentAuditor", "only auditor can cancel audit")
	}

	if d.isPosted {
		return commonErrors.NewSlugError("voucher-cancelAudit-posted", "voucher is posted")
	}

	d.isAudited = false
	d.auditor = uuid.Nil
	return nil
}
