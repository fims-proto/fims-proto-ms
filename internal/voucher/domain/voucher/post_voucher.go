package voucher

import (
	"github.com/google/uuid"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (d *Voucher) Post(poster uuid.UUID) error {
	if d.isPosted {
		return commonErrors.NewSlugError("voucher-post-repeatPost", "voucher is posted")
	}

	if !d.isAudited {
		return commonErrors.NewSlugError("voucher-post-notAudited", "voucher is not audited")
	}

	if !d.isReviewed {
		return commonErrors.NewSlugError("voucher-post-notReviewed", "voucher is not reviewed")
	}

	if poster == uuid.Nil {
		return commonErrors.NewSlugError("voucher-post-emptyPoster", "empty poster")
	}

	d.isPosted = true
	d.poster = poster
	return nil
}
