package report

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

// AddItemToSection adds a new item to the specified section at the given position
// insertAfterSequence:
//   - 0 or not provided: insert at the beginning (sequence=1)
//   - N: insert after sequence=N (new item becomes sequence=N+1)
//   - >= max sequence: append to the end
func (r *Report) AddItemToSection(sectionId uuid.UUID, item *Item, insertAfterSequence int) error {
	if item == nil {
		return errors.NewSlugError("report-item-nil")
	}

	section, err := r.findSectionById(sectionId)
	if err != nil {
		return err
	}

	return section.AddItem(item, insertAfterSequence)
}

// DeleteItemFromSection deletes an item from the specified section
func (r *Report) DeleteItemFromSection(sectionId uuid.UUID, itemId uuid.UUID) error {
	section, err := r.findSectionById(sectionId)
	if err != nil {
		return err
	}

	return section.DeleteItem(itemId)
}

// findSectionById recursively searches for a section by ID
func (r *Report) findSectionById(sectionId uuid.UUID) (*Section, error) {
	for _, section := range r.sections {
		if found := section.findSectionByIdRecursive(sectionId); found != nil {
			return found, nil
		}
	}
	return nil, errors.NewSlugError("report-section-notFound", map[string]interface{}{"sectionId": sectionId.String()})
}

// AddItem adds an item at the specified position and renumbers sequences
// insertAfterSequence:
//   - 0 or not provided: insert at the beginning (sequence=1)
//   - N: insert after sequence=N (new item becomes sequence=N+1)
//   - >= max sequence: append to the end
func (s *Section) AddItem(item *Item, insertAfterSequence int) error {
	if item == nil {
		return errors.NewSlugError("report-item-nil")
	}

	// If no items exist or insertAfterSequence <= 0, insert at beginning
	if len(s.items) == 0 || insertAfterSequence <= 0 {
		s.items = append([]*Item{item}, s.items...)
		s.renumberItems()
		return nil
	}

	// If insertAfterSequence >= current max, append to end
	if insertAfterSequence >= len(s.items) {
		s.items = append(s.items, item)
		s.renumberItems()
		return nil
	}

	// Insert at specified position (after insertAfterSequence)
	insertIndex := insertAfterSequence // Insert at position insertAfterSequence (0-based index)
	s.items = append(s.items[:insertIndex], append([]*Item{item}, s.items[insertIndex:]...)...)

	// Renumber all items to ensure sequential ordering
	s.renumberItems()

	return nil
}

// DeleteItem removes an item from this section by ID
func (s *Section) DeleteItem(itemId uuid.UUID) error {
	// Find the item and check if it's editable
	itemIndex := -1
	for i, item := range s.items {
		if item.id == itemId {
			if !item.isEditable {
				return errors.NewSlugError("report-item-notEditable")
			}
			itemIndex = i
			break
		}
	}

	if itemIndex == -1 {
		return errors.NewSlugError("report-item-notFound", map[string]interface{}{"itemId": itemId.String()})
	}

	// Remove the item from the slice
	s.items = append(s.items[:itemIndex], s.items[itemIndex+1:]...)

	// Renumber remaining items to maintain sequential sequences
	s.renumberItems()

	return nil
}

// renumberItems assigns sequential sequence numbers (1, 2, 3, ...) to all items
func (s *Section) renumberItems() {
	for i, item := range s.items {
		item.SetSequence(i + 1)
	}
}

// findSectionByIdRecursive searches this section and all nested sections for the specified ID
func (s *Section) findSectionByIdRecursive(sectionId uuid.UUID) *Section {
	if s.id == sectionId {
		return s
	}

	for _, childSection := range s.sections {
		if found := childSection.findSectionByIdRecursive(sectionId); found != nil {
			return found
		}
	}

	return nil
}
