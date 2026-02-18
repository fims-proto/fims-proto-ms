package report

import (
	"testing"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/formula_rule"

	"github.com/google/uuid"
)

func TestReport_UpdateReportStructure_UpdateTitle(t *testing.T) {
	// Setup: Create a report
	section1 := createTestSection("Section 1", 1, nil, nil)
	report := createTestReport([]*Section{section1})

	// Execute: Update title
	newTitle := "Updated Title"
	params := UpdateReportParams{
		Title:    &newTitle,
		Sections: []UpdateReportParamsSection{},
	}

	err := report.UpdateReportStructure(params)
	// Verify
	if err != nil {
		t.Errorf("UpdateReportStructure() error = %v, want nil", err)
	}
	if report.title != "Updated Title" {
		t.Errorf("Report title should be 'Updated Title', got %s", report.title)
	}
}

func TestReport_UpdateReportStructure_UpdateAmountTypes(t *testing.T) {
	// Setup
	section1 := createTestSection("Section 1", 1, nil, nil)
	report := createTestReport([]*Section{section1})

	// Execute: Update amount types
	periodAmountType, _ := amount_type.FromString("period_amount")
	params := UpdateReportParams{
		AmountTypes: []amount_type.AmountType{periodAmountType},
		Sections:    []UpdateReportParamsSection{},
	}

	err := report.UpdateReportStructure(params)
	// Verify
	if err != nil {
		t.Errorf("UpdateReportStructure() error = %v, want nil", err)
	}
	if len(report.amountTypes) != 1 {
		t.Errorf("Report should have 1 amount type, got %d", len(report.amountTypes))
	}
}

func TestSection_SynchronizeItems_AddNewItem(t *testing.T) {
	// Setup: Section with one item
	existingItem := createTestItem("Existing Item", 1, true)
	section := createTestSection("Test Section", 1, nil, []*Item{existingItem})
	existingItemId := existingItem.id

	// Execute: Add new item (no ID, sequence 2)
	newItemText := "New Item"
	newItemLevel := 1
	newItemSumFactor := 1
	newItemDataSource, _ := data_source.FromString("sum")

	params := []UpdateReportParamsItem{
		// Keep existing item
		{
			ItemId: &existingItemId,
		},
		// New item
		{
			ItemId:     nil, // No ID = create new
			Text:       &newItemText,
			Level:      &newItemLevel,
			SumFactor:  &newItemSumFactor,
			DataSource: &newItemDataSource,
		},
	}

	err := section.SynchronizeItems(params)
	// Verify
	if err != nil {
		t.Errorf("SynchronizeItems() error = %v, want nil", err)
	}
	if len(section.items) != 2 {
		t.Errorf("Section should have 2 items, got %d", len(section.items))
	}
	if section.items[1].text != "New Item" {
		t.Errorf("Second item text should be 'New Item', got %s", section.items[1].text)
	}
}

func TestSection_SynchronizeItems_DeleteItem(t *testing.T) {
	// Setup: Section with 2 items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2})
	item1Id := item1.id

	// Execute: Only include item1 (implicit delete of item2)
	params := []UpdateReportParamsItem{
		{
			ItemId: &item1Id,
		},
		// item2 not included = deleted
	}

	err := section.SynchronizeItems(params)
	// Verify
	if err != nil {
		t.Errorf("SynchronizeItems() error = %v, want nil", err)
	}
	if len(section.items) != 1 {
		t.Errorf("Section should have 1 item, got %d", len(section.items))
	}
	if section.items[0].id != item1Id {
		t.Error("Remaining item should be item1")
	}
}

func TestSection_SynchronizeItems_UpdateItemText(t *testing.T) {
	// Setup
	item1 := createTestItem("Original Text", 1, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1})
	itemId := item1.id

	// Execute: Update item text
	newText := "Updated Text"
	params := []UpdateReportParamsItem{
		{
			ItemId: &itemId,
			Text:   &newText,
		},
	}

	err := section.SynchronizeItems(params)
	// Verify
	if err != nil {
		t.Errorf("SynchronizeItems() error = %v, want nil", err)
	}
	if section.items[0].text != "Updated Text" {
		t.Errorf("Item text should be 'Updated Text', got %s", section.items[0].text)
	}
}

func TestSection_SynchronizeItems_ReorderItems(t *testing.T) {
	// Setup: 3 items in sequence
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	item3 := createTestItem("Item 3", 3, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2, item3})
	item1Id := item1.id
	item2Id := item2.id
	item3Id := item3.id

	// Execute: Reverse order (3, 2, 1)
	params := []UpdateReportParamsItem{
		{ItemId: &item3Id},
		{ItemId: &item2Id},
		{ItemId: &item1Id},
	}

	err := section.SynchronizeItems(params)
	// Verify
	if err != nil {
		t.Errorf("SynchronizeItems() error = %v, want nil", err)
	}
	if section.items[0].id != item3Id {
		t.Error("First item should be item3")
	}
	if section.items[1].id != item2Id {
		t.Error("Second item should be item2")
	}
	if section.items[2].id != item1Id {
		t.Error("Third item should be item1")
	}
	// Verify sequences were renumbered
	if section.items[0].sequence != 1 || section.items[1].sequence != 2 || section.items[2].sequence != 3 {
		t.Error("Items should have sequential sequences 1, 2, 3")
	}
}

func TestSection_SynchronizeItems_ComplexBatchOperation(t *testing.T) {
	// Setup: 2 items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2})
	item1Id := item1.id

	// Execute: Delete item2, add new item, update item1 text, reorder
	updatedText := "Updated Item 1"
	newItemText := "New Item"
	newItemLevel := 1
	newItemSumFactor := 1
	newItemDataSource, _ := data_source.FromString("sum")

	params := []UpdateReportParamsItem{
		// Updated item1 first
		{
			ItemId: &item1Id,
			Text:   &updatedText,
		},
		// New item second
		{
			ItemId:     nil,
			Text:       &newItemText,
			Level:      &newItemLevel,
			SumFactor:  &newItemSumFactor,
			DataSource: &newItemDataSource,
		},
		// item2 deleted (not included)
	}

	err := section.SynchronizeItems(params)
	// Verify
	if err != nil {
		t.Errorf("SynchronizeItems() error = %v, want nil", err)
	}
	if len(section.items) != 2 {
		t.Errorf("Section should have 2 total items, got %d", len(section.items))
	}
	if section.items[0].text != "Updated Item 1" {
		t.Errorf("First item should be 'Updated Item 1', got %s", section.items[0].text)
	}
	if section.items[1].text != "New Item" {
		t.Errorf("Second item should be 'New Item', got %s", section.items[1].text)
	}
}

func TestSection_SynchronizeItems_ErrorNonEditableItem(t *testing.T) {
	// Setup: Report with non-editable item
	nonEditableItem := createTestItem("Non-Editable Item", 1, false)
	section := createTestSection("Test Section", 1, nil, []*Item{nonEditableItem})

	// Execute: Try to delete non-editable item by not including it
	params := []UpdateReportParamsItem{} // Empty = delete all items

	err := section.SynchronizeItems(params)

	// Verify
	if err == nil {
		t.Error("SynchronizeItems() should return error for deleting non-editable item")
	}
}

func TestSection_SynchronizeItems_UpdateSumFactor(t *testing.T) {
	// Setup
	item1 := createTestItem("Item 1", 1, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1})
	itemId := item1.id

	// Execute: Update sum factor
	newSumFactor := -1
	params := []UpdateReportParamsItem{
		{
			ItemId:    &itemId,
			SumFactor: &newSumFactor,
		},
	}

	err := section.SynchronizeItems(params)
	// Verify
	if err != nil {
		t.Errorf("SynchronizeItems() error = %v, want nil", err)
	}
	if section.items[0].sumFactor != -1 {
		t.Errorf("Item sumFactor should be -1, got %d", section.items[0].sumFactor)
	}
}

func TestSection_SynchronizeItems_UpdateDataSourceAndFormulas(t *testing.T) {
	// Setup
	item1 := createTestItem("Item 1", 1, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1})
	itemId := item1.id

	// Execute: Update data source to formulas
	accountId := uuid.New()
	formulasDataSource, _ := data_source.FromString("formulas")
	formulaRule, _ := formula_rule.FromString("debit")

	params := []UpdateReportParamsItem{
		{
			ItemId:     &itemId,
			DataSource: &formulasDataSource,
			Formulas: []UpdateReportParamsFormula{
				{
					SumFactor: 1,
					AccountId: accountId,
					Rule:      formulaRule,
				},
			},
		},
	}

	err := section.SynchronizeItems(params)
	// Verify
	if err != nil {
		t.Errorf("SynchronizeItems() error = %v, want nil", err)
	}
	if section.items[0].dataSource != data_source.Formulas {
		t.Errorf("Item dataSource should be Formulas, got %v", section.items[0].dataSource)
	}
	if len(section.items[0].formulas) != 1 {
		t.Errorf("Item should have 1 formula, got %d", len(section.items[0].formulas))
	}
}

func TestSection_createItemFromData_MissingRequiredFields(t *testing.T) {
	section := createTestSection("Test Section", 1, nil, nil)

	// Test missing text
	params := UpdateReportParamsItem{
		Text:       nil,
		Level:      intPtr(1),
		SumFactor:  intPtr(1),
		DataSource: dataSourcePtr(data_source.Sum),
	}
	_, err := section.createItemFromData(params, 1)
	if err == nil {
		t.Error("createItemFromData() should error when text is missing")
	}

	// Test missing level
	params = UpdateReportParamsItem{
		Text:       stringPtr("Test"),
		Level:      nil,
		SumFactor:  intPtr(1),
		DataSource: dataSourcePtr(data_source.Sum),
	}
	_, err = section.createItemFromData(params, 1)
	if err == nil {
		t.Error("createItemFromData() should error when level is missing")
	}

	// Test missing sumFactor
	params = UpdateReportParamsItem{
		Text:       stringPtr("Test"),
		Level:      intPtr(1),
		SumFactor:  nil,
		DataSource: dataSourcePtr(data_source.Sum),
	}
	_, err = section.createItemFromData(params, 1)
	if err == nil {
		t.Error("createItemFromData() should error when sumFactor is missing")
	}

	// Test missing dataSource
	params = UpdateReportParamsItem{
		Text:       stringPtr("Test"),
		Level:      intPtr(1),
		SumFactor:  intPtr(1),
		DataSource: nil,
	}
	_, err = section.createItemFromData(params, 1)
	if err == nil {
		t.Error("createItemFromData() should error when dataSource is missing")
	}
}

func TestItem_applyUpdates_NonEditableItem(t *testing.T) {
	// Setup: Non-editable item
	item := createTestItem("Non-Editable Item", 1, false)

	// Execute: Try to update
	newText := "Updated Text"
	params := UpdateReportParamsItem{
		Text: &newText,
	}

	err := item.applyUpdates(params)

	// Verify
	if err == nil {
		t.Error("applyUpdates() should return error for non-editable item")
	}
}

func TestReport_UpdateReportStructure_UpdateSectionTitle(t *testing.T) {
	// Setup
	section1 := createTestSection("Original Title", 1, nil, nil)
	report := createTestReport([]*Section{section1})
	sectionId := section1.id

	// Execute: Update section title
	newSectionTitle := "Updated Section Title"
	params := UpdateReportParams{
		Sections: []UpdateReportParamsSection{
			{
				SectionId: sectionId,
				Title:     &newSectionTitle,
				Items:     []UpdateReportParamsItem{},
			},
		},
	}

	err := report.UpdateReportStructure(params)
	// Verify
	if err != nil {
		t.Errorf("UpdateReportStructure() error = %v, want nil", err)
	}
	section, _ := report.findSectionById(sectionId)
	if section.title != "Updated Section Title" {
		t.Errorf("Section title should be 'Updated Section Title', got %s", section.title)
	}
}

// Test that non-editable items can have their sequences updated during reordering
func TestSynchronizeItems_NonEditableItemSequenceUpdates(t *testing.T) {
	// Setup: Section with 3 items: editable (seq 1), non-editable (seq 2), editable (seq 3)
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2 - Non-Editable", 2, false) // Non-editable
	item3 := createTestItem("Item 3", 3, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2, item3})

	item1Id := item1.id
	item2Id := item2.id
	item3Id := item3.id

	// Execute: Reorder existing items and add a new item
	// Place existing items in order: item1, item3, item2 (non-editable moved to position 3)
	// Then add new item (will be appended as position 4)
	newItemText := "New Item - Inserted"
	newItemLevel := 1
	newItemSumFactor := 1
	newItemDataSource, _ := data_source.FromString("sum")

	params := []UpdateReportParamsItem{
		{ItemId: &item1Id},
		{ItemId: &item3Id}, // Item 3 moves to position 2
		{ItemId: &item2Id}, // Non-editable item moves to position 3
		// New item will be added after existing items
		{
			ItemId:     nil,
			Text:       &newItemText,
			Level:      &newItemLevel,
			SumFactor:  &newItemSumFactor,
			DataSource: &newItemDataSource,
		},
	}

	err := section.SynchronizeItems(params)
	// Verify
	if err != nil {
		t.Fatalf("SynchronizeItems() error = %v, want nil", err)
	}

	// Should have 4 items total
	if len(section.items) != 4 {
		t.Errorf("Section should have 4 items, got %d", len(section.items))
	}

	// Verify all items have sequential sequences (1, 2, 3, 4)
	for i, item := range section.items {
		expectedSeq := i + 1
		if item.sequence != expectedSeq {
			t.Errorf("Item %d should have sequence %d, got %d", i, expectedSeq, item.sequence)
		}
	}

	// Verify the ordering: item1, item3, item2 (non-editable), new item
	if section.items[0].id != item1Id {
		t.Error("Item 1 should be at position 1")
	}
	if section.items[1].id != item3Id {
		t.Error("Item 3 should be at position 2")
	}
	if section.items[2].id != item2Id {
		t.Error("Non-editable item should be at position 3")
	}
	if section.items[2].sequence != 3 {
		t.Errorf("Non-editable item should have sequence 3, got %d", section.items[2].sequence)
	}

	// Verify new item was created at position 4
	if section.items[3].text != "New Item - Inserted" {
		t.Errorf("New item should be at position 4, got %s at position 4", section.items[3].text)
	}
}

// Test updating report with nested sections
func TestUpdateReportStructure_WithNestedSections(t *testing.T) {
	// Setup: Create report with nested section structure
	// Section 1 (top-level) containing items A, B
	itemA := createTestItem("Item A", 1, true)
	itemB := createTestItem("Item B", 2, true)
	section1 := createTestSection("Section 1", 1, nil, []*Item{itemA, itemB})

	// Section 2 (top-level) with items C, D and Subsection 2.1 with items E, F
	itemC := createTestItem("Item C", 1, true)
	itemD := createTestItem("Item D", 2, true)
	itemE := createTestItem("Item E", 1, true)
	itemF := createTestItem("Item F", 2, true)
	subsection21 := createTestSection("Subsection 2.1", 1, nil, []*Item{itemE, itemF})
	section2 := createTestSection("Section 2", 2, []*Section{subsection21}, []*Item{itemC, itemD})

	report := createTestReport([]*Section{section1, section2})

	// Execute: Update report with:
	// - Section 1: update title, update item A's text
	// - Section 2: keep title, update item C, update subsection 2.1's title
	updatedSection1Title := "Section 1 - Updated"
	updatedItemAText := "Item A - Updated"
	updatedItemCText := "Item C - Updated"
	updatedSubsection21Title := "Subsection 2.1 - Updated"

	itemAId := itemA.id
	itemBId := itemB.id
	itemCId := itemC.id
	itemDId := itemD.id
	itemEId := itemE.id
	itemFId := itemF.id

	params := UpdateReportParams{
		Sections: []UpdateReportParamsSection{
			{
				SectionId: section1.id,
				Title:     &updatedSection1Title,
				Items: []UpdateReportParamsItem{
					{
						ItemId: &itemAId,
						Text:   &updatedItemAText,
					},
					{
						ItemId: &itemBId,
					},
				},
			},
			{
				SectionId: section2.id,
				Items: []UpdateReportParamsItem{
					{
						ItemId: &itemCId,
						Text:   &updatedItemCText,
					},
					{
						ItemId: &itemDId,
					},
				},
				Sections: []UpdateReportParamsSection{
					{
						SectionId: subsection21.id,
						Title:     &updatedSubsection21Title,
						Items: []UpdateReportParamsItem{
							{ItemId: &itemEId},
							{ItemId: &itemFId},
						},
					},
				},
			},
		},
	}

	err := report.UpdateReportStructure(params)
	// Verify
	if err != nil {
		t.Fatalf("UpdateReportStructure() error = %v, want nil", err)
	}

	// Verify section titles updated correctly
	if report.sections[0].title != "Section 1 - Updated" {
		t.Errorf("Section 1 title should be 'Section 1 - Updated', got %s", report.sections[0].title)
	}
	if report.sections[1].title != "Section 2" {
		t.Errorf("Section 2 title should remain 'Section 2', got %s", report.sections[1].title)
	}

	// Verify items updated correctly at all levels
	if report.sections[0].items[0].text != "Item A - Updated" {
		t.Errorf("Item A text should be 'Item A - Updated', got %s", report.sections[0].items[0].text)
	}
	if report.sections[1].items[0].text != "Item C - Updated" {
		t.Errorf("Item C text should be 'Item C - Updated', got %s", report.sections[1].items[0].text)
	}

	// Verify subsection title updated
	if report.sections[1].sections[0].title != "Subsection 2.1 - Updated" {
		t.Errorf("Subsection 2.1 title should be 'Subsection 2.1 - Updated', got %s", report.sections[1].sections[0].title)
	}

	// Verify nested section structure preserved
	if len(report.sections[1].sections) != 1 {
		t.Errorf("Section 2 should have 1 subsection, got %d", len(report.sections[1].sections))
	}

	// Verify sequences recalculated for all levels
	for i, item := range report.sections[0].items {
		if item.sequence != i+1 {
			t.Errorf("Section 1 item %d should have sequence %d, got %d", i, i+1, item.sequence)
		}
	}
	for i, item := range report.sections[1].items {
		if item.sequence != i+1 {
			t.Errorf("Section 2 item %d should have sequence %d, got %d", i, i+1, item.sequence)
		}
	}
	for i, item := range report.sections[1].sections[0].items {
		if item.sequence != i+1 {
			t.Errorf("Subsection 2.1 item %d should have sequence %d, got %d", i, i+1, item.sequence)
		}
	}
}

// Helper functions for tests

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func dataSourcePtr(ds data_source.DataSource) *data_source.DataSource {
	return &ds
}
