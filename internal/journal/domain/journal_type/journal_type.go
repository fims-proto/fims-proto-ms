package journal_type

import (
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

type JournalType struct {
	slug string
}

func (t JournalType) String() string {
	return t.slug
}

var (
	Unknown = JournalType{""}
	General = JournalType{"general_journal"}
)

func FromString(s string) (JournalType, error) {
	switch s {
	case General.slug:
		return General, nil
	}

	return Unknown, errors.NewSlugError("journalEntry-unknownJournalType", fmt.Sprintf("unknown journal type %s", s), s)
}
