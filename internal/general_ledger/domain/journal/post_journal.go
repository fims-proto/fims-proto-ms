package journal

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (j *Journal) Post(poster string) error {
	if j.Period().IsClosed() {
		return errors.NewInvalidInputError(errors.SlugJournalPostPeriodClosed)
	}

	if !j.Period().IsCurrent() {
		return errors.NewInvalidInputError(errors.SlugJournalPostPeriodNotCurrent)
	}

	if j.isPosted {
		return errors.NewInvalidInputError(errors.SlugJournalPostRepeatPost)
	}

	if !j.isAudited {
		return errors.NewInvalidInputError(errors.SlugJournalPostNotAudited)
	}

	if !j.isReviewed {
		return errors.NewInvalidInputError(errors.SlugJournalPostNotReviewed)
	}

	if isEmptyUser(poster) {
		return errors.NewInternalError(errors.SlugJournalPostEmptyPoster)
	}

	j.isPosted = true
	j.poster = poster
	return nil
}
