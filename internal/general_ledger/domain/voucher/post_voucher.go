package voucher

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

func (v *Voucher) Post(poster uuid.UUID) error {
	if v.Period().IsClosed() {
		return errors.NewSlugError("voucher-post-periodClosed")
	}

	if !v.Period().IsCurrent() {
		return errors.NewSlugError("voucher-post-periodNotCurrent")
	}

	if v.isPosted {
		return errors.NewSlugError("voucher-post-repeatPost")
	}

	if !v.isAudited {
		return errors.NewSlugError("voucher-post-notAudited")
	}

	if !v.isReviewed {
		return errors.NewSlugError("voucher-post-notReviewed")
	}

	if poster == uuid.Nil {
		return errors.NewSlugError("voucher-post-emptyPoster")
	}

	v.isPosted = true
	v.poster = poster
	return nil
}
