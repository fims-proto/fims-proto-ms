package journal

type JournalType string

const (
	TypeGeneral       JournalType = "GENERAL"
	TypeAdjusting     JournalType = "ADJUSTING"
	TypeReversing     JournalType = "REVERSING"
	TypeClosing       JournalType = "CLOSING"
	TypeYearlyClosing JournalType = "YEARLY_CLOSING"
)

func (t JournalType) IsValid() bool {
	switch t {
	case TypeGeneral, TypeAdjusting, TypeReversing, TypeClosing, TypeYearlyClosing:
		return true
	}
	return false
}

func (t JournalType) RequiresReferenceJournal() bool {
	return t == TypeAdjusting || t == TypeReversing
}
