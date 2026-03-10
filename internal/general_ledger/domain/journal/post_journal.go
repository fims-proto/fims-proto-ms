package journal

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

func (j *Journal) Post(poster uuid.UUID) error {
	if j.Period().IsClosed() {
		return errors.NewSlugError("journal-post-periodClosed")
	}

	if !j.Period().IsCurrent() {
		return errors.NewSlugError("journal-post-periodNotCurrent")
	}

	if j.isPosted {
		return errors.NewSlugError("journal-post-repeatPost")
	}

	if !j.isAudited {
		return errors.NewSlugError("journal-post-notAudited")
	}

	if !j.isReviewed {
		return errors.NewSlugError("journal-post-notReviewed")
	}

	if poster == uuid.Nil {
		return errors.NewSlugError("journal-post-emptyPoster")
	}

	j.isPosted = true
	j.poster = poster
	return nil
}
