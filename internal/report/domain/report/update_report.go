package report

import (
	"sort"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/formula_rule"

	"github.com/google/uuid"
)

// UpdateReportParams contains parameters for updating a report
type UpdateReportParams struct {
	Title       *string
	AmountTypes []amount_type.AmountType
	Sections    []UpdateReportParamsSection
}

// UpdateReportParamsSection contains parameters for updating a section
type UpdateReportParamsSection struct {
	SectionId uuid.UUID
	Title     *string
	Items     []UpdateReportParamsItem
	Sections  []UpdateReportParamsSection
}

// UpdateReportParamsItem contains parameters for updating or creating an item
type UpdateReportParamsItem struct {
	ItemId           *uuid.UUID
	Text             *string
	Level            *int
	SumFactor        *int
	DisplaySumFactor *bool
	DataSource       *data_source.DataSource
	Formulas         []UpdateReportParamsFormula
	IsBreakdownItem  *bool
	IsAbleToAddChild *bool
}

// UpdateReportParamsFormula contains parameters for a formula
type UpdateReportParamsFormula struct {
	FormulaId *uuid.UUID
	SumFactor int
	AccountId uuid.UUID
	Rule      formula_rule.FormulaRule
}

// UpdateReportStructure applies comprehensive updates: report metadata, section updates, and item CRUD operations
func (r *Report) UpdateReportStructure(params UpdateReportParams) error {
	// 1. Update report-level fields
	if params.Title != nil {
		r.title = *params.Title
	}

	if len(params.AmountTypes) > 0 {
		r.amountTypes = params.AmountTypes
	}

	// 2. Update each section
	for _, sectionData := range params.Sections {
		section, err := r.findSectionById(sectionData.SectionId)
		if err != nil {
			return err
		}

		// Update section title if provided
		if sectionData.Title != nil {
			section.title = *sectionData.Title
		}

		// Apply item updates (add, update, delete, reorder)
		err = section.SynchronizeItems(sectionData.Items)
		if err != nil {
			return err
		}

		// Apply subsection updates (recursive)
		if len(sectionData.Sections) > 0 {
			if err = section.SynchronizeSubsections(sectionData.Sections); err != nil {
				return err
			}
		}
	}

	return nil
}

// SynchronizeItems performs a diff between current items and desired items, then:
// - Creates new items (items without ID)
// - Updates existing items (items with ID and update fields)
// - Deletes missing items (current items not in desired list)
// - Reorders all items by sequence
func (s *Section) SynchronizeItems(desiredItems []UpdateReportParamsItem) error {
	// Build map of desired items by ID
	desiredById := make(map[uuid.UUID]*UpdateReportParamsItem)
	var newItems []*UpdateReportParamsItem

	for i := range desiredItems {
		item := &desiredItems[i]
		if item.ItemId == nil {
			// New item to create
			newItems = append(newItems, item)
		} else {
			// Existing item
			desiredById[*item.ItemId] = item
		}
	}

	// 1. Update existing items and identify items to delete
	var keptItems []*Item
	for _, currentItem := range s.items {
		if desiredItem, exists := desiredById[currentItem.id]; exists {
			// Item still exists - apply updates
			if err := currentItem.applyUpdates(*desiredItem); err != nil {
				return err
			}
			keptItems = append(keptItems, currentItem)
		} else {
			// Item not in desired list - delete it
			if !currentItem.isEditable {
				return errors.NewSlugError("report-item-notEditable")
			}
			// Simply don't add to keptItems (implicit deletion)
		}
	}

	// 2. Create new items
	for i, newItemData := range newItems {
		// Sequence is determined by position: len(keptItems) + i + 1
		newItem, err := s.createItemFromData(*newItemData, len(keptItems)+i+1)
		if err != nil {
			return err
		}

		keptItems = append(keptItems, newItem)
	}

	// 3. Replace section's items
	s.items = keptItems

	// 4. Reorder items by sequence (using array position)
	if err := s.reorderItemsBySequence(desiredItems); err != nil {
		return err
	}

	return nil
}

// applyUpdates updates item fields if provided in the update data
func (i *Item) applyUpdates(data UpdateReportParamsItem) error {
	if data.Text != nil {
		if err := i.UpdateText(*data.Text); err != nil {
			return err
		}
	}

	if data.SumFactor != nil {
		if err := i.UpdateSumFactor(*data.SumFactor); err != nil {
			return err
		}
	}

	if data.DataSource != nil {
		// Build formulas
		var formulas []*Formula
		for _, fData := range data.Formulas {
			formulaId := uuid.New()
			if fData.FormulaId != nil {
				formulaId = *fData.FormulaId
			}
			formula, err := NewFormula(formulaId, len(formulas)+1, fData.AccountId, fData.SumFactor, fData.Rule.String(), nil)
			if err != nil {
				return err
			}
			formulas = append(formulas, formula)
		}

		if err := i.UpdateDataSource(*data.DataSource, formulas); err != nil {
			return err
		}
	}

	// Note: Level, DisplaySumFactor, ItemType, IsBreakdownItem, IsAbleToAddChild
	// are typically not updated after creation, but can be added if needed

	return nil
}

// createItemFromData creates a new item from update data
func (s *Section) createItemFromData(data UpdateReportParamsItem, sequence int) (*Item, error) {
	// Validate required fields for new items
	if data.Text == nil {
		return nil, errors.NewSlugError("report-item-textRequired")
	}
	if data.Level == nil {
		return nil, errors.NewSlugError("report-item-levelRequired")
	}
	if data.SumFactor == nil {
		return nil, errors.NewSlugError("report-item-sumFactorRequired")
	}
	if data.DataSource == nil {
		return nil, errors.NewSlugError("report-item-dataSourceRequired")
	}

	// Build formulas
	var formulas []*Formula
	for _, fData := range data.Formulas {
		formula, err := NewFormula(uuid.New(), len(formulas)+1, fData.AccountId, fData.SumFactor, fData.Rule.String(), nil)
		if err != nil {
			return nil, err
		}
		formulas = append(formulas, formula)
	}

	// Default values
	displaySumFactor := false
	if data.DisplaySumFactor != nil {
		displaySumFactor = *data.DisplaySumFactor
	}

	// ItemType is not provided in update request - use empty string for new items
	itemTypeStr := ""

	isBreakdownItem := false
	if data.IsBreakdownItem != nil {
		isBreakdownItem = *data.IsBreakdownItem
	}

	isAbleToAddChild := false
	if data.IsAbleToAddChild != nil {
		isAbleToAddChild = *data.IsAbleToAddChild
	}

	// Create new item
	item, err := NewItem(
		uuid.New(),
		*data.Text,
		*data.Level,
		sequence, // Use provided sequence
		itemTypeStr,
		*data.SumFactor,
		displaySumFactor,
		data.DataSource.String(),
		formulas,
		nil,  // amounts - will be calculated during report generation
		true, // isEditable - new items are editable
		isBreakdownItem,
		isAbleToAddChild,
	)
	if err != nil {
		return nil, err
	}

	return item, nil
}

// reorderItemsBySequence reorders items to match the order in desiredItems array
func (s *Section) reorderItemsBySequence(desiredItems []UpdateReportParamsItem) error {
	// Step 1: Build a set of IDs that appear in desiredItems with non-nil ItemId
	// This allows us to classify items in s.items as existing vs new
	desiredIDs := make(map[uuid.UUID]bool)
	for _, desired := range desiredItems {
		if desired.ItemId != nil {
			desiredIDs[*desired.ItemId] = true
		}
	}

	// Step 2: Classify items in s.items as existing (by ID) vs new
	itemsByID := make(map[uuid.UUID]*Item) // Existing items
	var newItems []*Item                   // New items (in order they appear in s.items)

	for _, item := range s.items {
		if desiredIDs[item.id] {
			// Existing item - appears in desiredItems with an ID
			itemsByID[item.id] = item
		} else {
			// New item - does not appear in desiredItems with an ID
			newItems = append(newItems, item)
		}
	}

	// Step 3: Build ordered list by following desiredItems array
	reorderedItems := make([]*Item, 0, len(s.items))
	newItemIdx := 0

	for _, desired := range desiredItems {
		if desired.ItemId != nil {
			// Existing item - look up by ID
			if item, ok := itemsByID[*desired.ItemId]; ok {
				reorderedItems = append(reorderedItems, item)
			}
		} else {
			// New item - take next from newItems list
			if newItemIdx < len(newItems) {
				reorderedItems = append(reorderedItems, newItems[newItemIdx])
				newItemIdx++
			}
		}
	}

	// Step 4: Replace section's items with reordered list
	s.items = reorderedItems

	// Step 5: Renumber all items sequentially (1, 2, 3, ...)
	s.renumberItems()

	return nil
}

// SynchronizeSubsections updates nested subsections recursively
// Only updates existing sections - does not create new sections
func (s *Section) SynchronizeSubsections(desiredSections []UpdateReportParamsSection) error {
	// Build map of desired sections by ID
	desiredById := make(map[uuid.UUID]*UpdateReportParamsSection)
	for i := range desiredSections {
		section := &desiredSections[i]
		desiredById[section.SectionId] = section
	}

	// Update existing subsections
	for _, currentSection := range s.sections {
		if desiredSection, exists := desiredById[currentSection.id]; exists {
			// Update section title if provided
			if desiredSection.Title != nil {
				currentSection.title = *desiredSection.Title
			}

			// Recursively update items
			if err := currentSection.SynchronizeItems(desiredSection.Items); err != nil {
				return err
			}

			// Recursively update nested subsections
			if len(desiredSection.Sections) > 0 {
				if err := currentSection.SynchronizeSubsections(desiredSection.Sections); err != nil {
					return err
				}
			}
		}
		// Note: We don't delete sections that aren't in the desired list
		// Only update existing ones
	}

	// Renumber subsections based on their position in desired array
	s.renumberSections(desiredSections)

	return nil
}

// renumberSections assigns sequential sequence numbers to subsections based on desired order
func (s *Section) renumberSections(desiredSections []UpdateReportParamsSection) {
	// Build map of section ID -> desired position
	positionMap := make(map[uuid.UUID]int)
	for i, desired := range desiredSections {
		positionMap[desired.SectionId] = i + 1
	}

	// Sort subsections by their desired position
	sort.SliceStable(s.sections, func(i, j int) bool {
		posI, hasI := positionMap[s.sections[i].id]
		posJ, hasJ := positionMap[s.sections[j].id]

		if hasI && hasJ {
			return posI < posJ
		}
		if hasI {
			return posI < s.sections[j].sequence
		}
		if hasJ {
			return s.sections[i].sequence < posJ
		}
		return s.sections[i].sequence < s.sections[j].sequence
	})

	// Assign sequential numbers
	for i, section := range s.sections {
		section.setSequence(i + 1)
	}
}

// setSequence sets the sequence of a section
func (s *Section) setSequence(sequence int) {
	s.sequence = sequence
}
