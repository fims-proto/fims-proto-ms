package report

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	commonerrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

type Section struct {
	id       uuid.UUID
	title    string
	sequence int // sequence within the parent, starts from 1
	amounts  []decimal.Decimal
	sections []*Section
	items    []*Item
}

func NewSection(
	id uuid.UUID,
	title string,
	sequence int,
	amounts []decimal.Decimal,
	sections []*Section,
	items []*Item,
) (*Section, error) {
	if id == uuid.Nil {
		return nil, errors.New("section id cannot be nil")
	}

	if sequence == 0 {
		return nil, commonerrors.NewSlugError("report-section-zeroSequence")
	}

	if len(sections) == 0 && len(items) == 0 {
		return nil, commonerrors.NewSlugError("report-section-emptySectionsAndItems")
	}

	return &Section{
		id:       id,
		title:    title,
		sequence: sequence,
		amounts:  amounts,
		sections: sections,
		items:    items,
	}, nil
}

func (s *Section) copy() *Section {
	var newSections []*Section
	for _, section := range s.sections {
		newSections = append(newSections, section.copy())
	}

	var newItems []*Item
	for _, item := range s.items {
		newItems = append(newItems, item.copy())
	}

	newSection, _ := NewSection(
		uuid.New(),
		s.title,
		s.sequence,
		nil,
		newSections,
		newItems,
	)
	return newSection
}

func (s *Section) SetAmounts(amounts []decimal.Decimal) {
	s.amounts = amounts
}

func (s *Section) Id() uuid.UUID {
	return s.id
}

func (s *Section) Title() string {
	return s.title
}

func (s *Section) Sequence() int {
	return s.sequence
}

func (s *Section) Amounts() []decimal.Decimal {
	return s.amounts
}

func (s *Section) Sections() []*Section {
	return s.sections
}

func (s *Section) Items() []*Item {
	return s.items
}
