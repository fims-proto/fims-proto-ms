package voucher

import (
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (d *Voucher) Audit(auditor uuid.UUID) error {
	if d.isAudited {
		return errors.NewSlugError("voucher-audit-repeatAudit")
	}

	if auditor == uuid.Nil {
		return errors.NewSlugError("voucher-audit-emptyAuditor")
	}

	if auditor == d.creator {
		return errors.NewSlugError("voucher-audit-auditorSameAsCreator")
	}

	if d.reviewer != uuid.Nil && auditor == d.reviewer {
		return errors.NewSlugError("voucher-audit-auditorSameAsReviewer")
	}

	d.isAudited = true
	d.auditor = auditor
	return nil
}

func (d *Voucher) CancelAudit(auditor uuid.UUID) error {
	if !d.isAudited {
		return errors.NewSlugError("voucher-cancelAudit-notAudited")
	}

	if d.auditor != auditor {
		return errors.NewSlugError("voucher-cancelAudit-differentAuditor")
	}

	if d.isPosted {
		return errors.NewSlugError("voucher-cancelAudit-posted")
	}

	d.isAudited = false
	d.auditor = uuid.Nil
	return nil
}
