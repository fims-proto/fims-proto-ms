package section_type

import "fmt"

type SectionType struct {
	slug string
}

func (st SectionType) String() string {
	return st.slug
}

var (
	None        = SectionType{""}
	Assets      = SectionType{"assets"}
	Liabilities = SectionType{"liabilities"}
	Equity      = SectionType{"equity"}
	Revenue     = SectionType{"revenue"}
	Expenses    = SectionType{"expenses"}
)

var stringToSectionType = map[string]SectionType{
	"":            None,
	"assets":      Assets,
	"liabilities": Liabilities,
	"equity":      Equity,
	"revenue":     Revenue,
	"expenses":    Expenses,
}

func FromString(s string) (SectionType, error) {
	sectionType, ok := stringToSectionType[s]
	if ok {
		return sectionType, nil
	}
	return None, fmt.Errorf("unknown section type: %s", s)
}
