package report

import (
	"fmt"
	"maps"
	"sort"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/formula_rule"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/item_type"

	"github.com/google/uuid"
)

// UpdateReportParams contains parameters for updating a report
type UpdateReportParams struct {
	Title       *string
	AmountTypes []amount_type.AmountType
	Sections    []UpdateSectionParams
}

// UpdateSectionParams contains parameters for updating a section
type UpdateSectionParams struct {
	SectionId uuid.UUID
	Title     *string
	Items     []UpdateItemParams
}

// UpdateItemParams contains parameters for updating or creating an item
type UpdateItemParams struct {
	ItemId           *uuid.UUID
	Sequence         int
	Text             *string
	Level            *int
	SumFactor        *int
	DisplaySumFactor *bool
	ItemType         *item_type.ItemType
	DataSource       *data_source.DataSource
	Formulas         []UpdateFormulaParams
	IsBreakdownItem  *bool
	IsAbleToAddChild *bool
}

// UpdateFormulaParams contains parameters for a formula
type UpdateFormulaParams struct {
	SumFactor int
	AccountId uuid.UUID
	Rule      formula_rule.FormulaRule
}

// UpdateReportStructure applies comprehensive updates: report metadata, section updates, and item CRUD operations
// Returns map of temporary client IDs -> actual UUIDs for newly created items
func (r *Report) UpdateReportStructure(params UpdateReportParams) (map[string]string, error) {
	createdItemIds := make(map[string]string)

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
			return nil, err
		}

		// Update section title if provided
		if sectionData.Title != nil {
			section.title = *sectionData.Title
		}

		// Apply item updates (add, update, delete, reorder)
		ids, err := section.SynchronizeItems(sectionData.Items)
		if err != nil {
			return nil, err
		}

		// Merge created item IDs
		maps.Copy(createdItemIds, ids)
	}

	return createdItemIds, nil
}

// SynchronizeItems performs a diff between current items and desired items, then:
// - Creates new items (items without ID)
// - Updates existing items (items with ID and update fields)
// - Deletes missing items (current items not in desired list)
// - Reorders all items by sequence
//
// Returns map of client temp IDs -> actual UUIDs for newly created items
func (s *Section) SynchronizeItems(desiredItems []UpdateItemParams) (map[string]string, error) {
	createdItemIds := make(map[string]string)

	// Build map of desired items by ID
	desiredById := make(map[uuid.UUID]*UpdateItemParams)
	var newItems []*UpdateItemParams

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
				return nil, err
			}
			keptItems = append(keptItems, currentItem)
		} else {
			// Item not in desired list - delete it
			if !currentItem.isEditable {
				return nil, errors.NewSlugError("report-item-notEditable")
			}
			// Simply don't add to keptItems (implicit deletion)
		}
	}

	// 2. Create new items
	for _, newItemData := range newItems {
		newItem, err := s.createItemFromData(*newItemData)
		if err != nil {
			return nil, err
		}

		keptItems = append(keptItems, newItem)

		// Track created item ID
		// Use sequence as temporary ID (frontend can send this to correlate)
		tempId := fmt.Sprintf("new-%d", newItemData.Sequence)
		createdItemIds[tempId] = newItem.id.String()
	}

	// 3. Replace section's items
	s.items = keptItems

	// 4. Reorder items by sequence
	if err := s.reorderItemsBySequence(desiredItems); err != nil {
		return nil, err
	}

	return createdItemIds, nil
}

// applyUpdates updates item fields if provided in the update data
func (i *Item) applyUpdates(data UpdateItemParams) error {
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
			formula, err := NewFormula(uuid.New(), len(formulas)+1, fData.AccountId, fData.SumFactor, fData.Rule.String(), nil)
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
func (s *Section) createItemFromData(data UpdateItemParams) (*Item, error) {
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

	itemTypeStr := ""
	if data.ItemType != nil {
		itemTypeStr = data.ItemType.String()
	}

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
		data.Sequence, // Will be renumbered later
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

// reorderItemsBySequence sorts items according to sequence in desired data
func (s *Section) reorderItemsBySequence(desiredItems []UpdateItemParams) error {
	// Build map of item ID -> desired sequence (for existing items)
	// And track all sequences to ensure new items are also ordered correctly
	sequenceMap := make(map[uuid.UUID]int)
	for _, desired := range desiredItems {
		if desired.ItemId != nil {
			sequenceMap[*desired.ItemId] = desired.Sequence
		}
	}

	// Sort items by their current sequence first (which was set during creation)
	// This ensures new items that have their sequence already set are ordered correctly
	sort.SliceStable(s.items, func(i, j int) bool {
		seqI, hasI := sequenceMap[s.items[i].id]
		seqJ, hasJ := sequenceMap[s.items[j].id]

		// Both have explicit sequences from the map
		if hasI && hasJ {
			return seqI < seqJ
		}
		// Only i has explicit sequence
		if hasI {
			// Compare with j's current sequence
			return seqI < s.items[j].sequence
		}
		// Only j has explicit sequence
		if hasJ {
			// Compare i's current sequence with j's explicit sequence
			return s.items[i].sequence < seqJ
		}

		// Neither has explicit sequence - use their current sequences
		return s.items[i].sequence < s.items[j].sequence
	})

	// Renumber all items sequentially (1, 2, 3, ...)
	s.renumberItems()

	return nil
}
