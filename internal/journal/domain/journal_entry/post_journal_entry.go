package journal_entry

import (
	"github.com/google/uuid"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (d *JournalEntry) Post(poster uuid.UUID) error {
	if d.isPosted {
		return commonErrors.NewSlugError("journalEntry-post-repeatPost", "entry is posted")
	}

	if !d.isAudited {
		return commonErrors.NewSlugError("journalEntry-post-notAudited", "entry is not audited")
	}

	if !d.isReviewed {
		return commonErrors.NewSlugError("journalEntry-post-notReviewed", "entry is not reviewed")
	}

	if poster == uuid.Nil {
		return commonErrors.NewSlugError("journalEntry-post-emptyPoster", "empty poster")
	}

	d.isPosted = true
	d.poster = poster
	return nil
}
