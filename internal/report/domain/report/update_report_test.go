package report

import (
	"testing"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/formula_rule"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/item_type"

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
		Sections: []UpdateSectionParams{},
	}

	_, err := report.UpdateReportStructure(params)

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
		Sections:    []UpdateSectionParams{},
	}

	_, err := report.UpdateReportStructure(params)

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

	params := []UpdateItemParams{
		// Keep existing item
		{
			ItemId:   &existingItemId,
			Sequence: 1,
		},
		// New item
		{
			ItemId:     nil, // No ID = create new
			Sequence:   2,
			Text:       &newItemText,
			Level:      &newItemLevel,
			SumFactor:  &newItemSumFactor,
			DataSource: &newItemDataSource,
		},
	}

	createdIds, err := section.SynchronizeItems(params)

	// Verify
	if err != nil {
		t.Errorf("SynchronizeItems() error = %v, want nil", err)
	}
	if len(createdIds) != 1 {
		t.Errorf("Should have created 1 item, got %d", len(createdIds))
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
	params := []UpdateItemParams{
		{
			ItemId:   &item1Id,
			Sequence: 1,
		},
		// item2 not included = deleted
	}

	_, err := section.SynchronizeItems(params)

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
	params := []UpdateItemParams{
		{
			ItemId:   &itemId,
			Sequence: 1,
			Text:     &newText,
		},
	}

	_, err := section.SynchronizeItems(params)

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
	params := []UpdateItemParams{
		{ItemId: &item3Id, Sequence: 1},
		{ItemId: &item2Id, Sequence: 2},
		{ItemId: &item1Id, Sequence: 3},
	}

	_, err := section.SynchronizeItems(params)

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

	params := []UpdateItemParams{
		// New item first
		{
			ItemId:     nil,
			Sequence:   1,
			Text:       &newItemText,
			Level:      &newItemLevel,
			SumFactor:  &newItemSumFactor,
			DataSource: &newItemDataSource,
		},
		// Updated item1 second
		{
			ItemId:   &item1Id,
			Sequence: 2,
			Text:     &updatedText,
		},
		// item2 deleted (not included)
	}

	createdIds, err := section.SynchronizeItems(params)

	// Verify
	if err != nil {
		t.Errorf("SynchronizeItems() error = %v, want nil", err)
	}
	if len(createdIds) != 1 {
		t.Errorf("Should have created 1 new item, got %d", len(createdIds))
	}
	if len(section.items) != 2 {
		t.Errorf("Section should have 2 total items, got %d", len(section.items))
	}
	if section.items[0].text != "New Item" {
		t.Errorf("First item should be 'New Item', got %s", section.items[0].text)
	}
	if section.items[1].text != "Updated Item 1" {
		t.Errorf("Second item should be 'Updated Item 1', got %s", section.items[1].text)
	}
}

func TestSection_SynchronizeItems_ErrorNonEditableItem(t *testing.T) {
	// Setup: Report with non-editable item
	nonEditableItem := createTestItem("Non-Editable Item", 1, false)
	section := createTestSection("Test Section", 1, nil, []*Item{nonEditableItem})

	// Execute: Try to delete non-editable item by not including it
	params := []UpdateItemParams{} // Empty = delete all items

	_, err := section.SynchronizeItems(params)

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
	params := []UpdateItemParams{
		{
			ItemId:    &itemId,
			Sequence:  1,
			SumFactor: &newSumFactor,
		},
	}

	_, err := section.SynchronizeItems(params)

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

	params := []UpdateItemParams{
		{
			ItemId:     &itemId,
			Sequence:   1,
			DataSource: &formulasDataSource,
			Formulas: []UpdateFormulaParams{
				{
					SumFactor: 1,
					AccountId: accountId,
					Rule:      formulaRule,
				},
			},
		},
	}

	_, err := section.SynchronizeItems(params)

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
	params := UpdateItemParams{
		Text:       nil,
		Level:      intPtr(1),
		SumFactor:  intPtr(1),
		DataSource: dataSourcePtr(data_source.Sum),
	}
	_, err := section.createItemFromData(params)
	if err == nil {
		t.Error("createItemFromData() should error when text is missing")
	}

	// Test missing level
	params = UpdateItemParams{
		Text:       stringPtr("Test"),
		Level:      nil,
		SumFactor:  intPtr(1),
		DataSource: dataSourcePtr(data_source.Sum),
	}
	_, err = section.createItemFromData(params)
	if err == nil {
		t.Error("createItemFromData() should error when level is missing")
	}

	// Test missing sumFactor
	params = UpdateItemParams{
		Text:       stringPtr("Test"),
		Level:      intPtr(1),
		SumFactor:  nil,
		DataSource: dataSourcePtr(data_source.Sum),
	}
	_, err = section.createItemFromData(params)
	if err == nil {
		t.Error("createItemFromData() should error when sumFactor is missing")
	}

	// Test missing dataSource
	params = UpdateItemParams{
		Text:       stringPtr("Test"),
		Level:      intPtr(1),
		SumFactor:  intPtr(1),
		DataSource: nil,
	}
	_, err = section.createItemFromData(params)
	if err == nil {
		t.Error("createItemFromData() should error when dataSource is missing")
	}
}

func TestItem_applyUpdates_NonEditableItem(t *testing.T) {
	// Setup: Non-editable item
	item := createTestItem("Non-Editable Item", 1, false)

	// Execute: Try to update
	newText := "Updated Text"
	params := UpdateItemParams{
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
		Sections: []UpdateSectionParams{
			{
				SectionId: sectionId,
				Title:     &newSectionTitle,
				Items:     []UpdateItemParams{},
			},
		},
	}

	_, err := report.UpdateReportStructure(params)

	// Verify
	if err != nil {
		t.Errorf("UpdateReportStructure() error = %v, want nil", err)
	}
	section, _ := report.findSectionById(sectionId)
	if section.title != "Updated Section Title" {
		t.Errorf("Section title should be 'Updated Section Title', got %s", section.title)
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

func itemTypePtr(it item_type.ItemType) *item_type.ItemType {
	return &it
}
