package voucher

import (
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (v *Voucher) Audit(auditor uuid.UUID) error {
	if v.isAudited {
		return errors.NewSlugError("voucher-audit-repeatAudit")
	}

	if auditor == uuid.Nil {
		return errors.NewSlugError("voucher-audit-emptyAuditor")
	}

	if auditor == v.creator {
		return errors.NewSlugError("voucher-audit-auditorSameAsCreator")
	}

	if v.reviewer != uuid.Nil && auditor == v.reviewer {
		return errors.NewSlugError("voucher-audit-auditorSameAsReviewer")
	}

	v.isAudited = true
	v.auditor = auditor
	return nil
}

func (v *Voucher) CancelAudit(auditor uuid.UUID) error {
	if !v.isAudited {
		return errors.NewSlugError("voucher-cancelAudit-notAudited")
	}

	if v.auditor != auditor {
		return errors.NewSlugError("voucher-cancelAudit-differentAuditor")
	}

	if v.isPosted {
		return errors.NewSlugError("voucher-cancelAudit-posted")
	}

	v.isAudited = false
	v.auditor = uuid.Nil
	return nil
}
