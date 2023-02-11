package voucher

import (
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (d *Voucher) Post(poster uuid.UUID) error {
	if d.isPosted {
		return errors.NewSlugError("voucher-post-repeatPost")
	}

	if !d.isAudited {
		return errors.NewSlugError("voucher-post-notAudited")
	}

	if !d.isReviewed {
		return errors.NewSlugError("voucher-post-notReviewed")
	}

	if poster == uuid.Nil {
		return errors.NewSlugError("voucher-post-emptyPoster")
	}

	d.isPosted = true
	d.poster = poster
	return nil
}
