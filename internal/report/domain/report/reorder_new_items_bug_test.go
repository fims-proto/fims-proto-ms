package report

import (
	"testing"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"

	"github.com/google/uuid"
)

// TestSection_SynchronizeItems_AddNewItemInMiddle tests the scenario where a new item
// is inserted in the middle of existing items. This test demonstrates the bug in
// reorderItemsBySequence function.
func TestSection_SynchronizeItems_AddNewItemInMiddle(t *testing.T) {
	// Setup: Section with 2 existing items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2})

	item1Id := item1.id
	item2Id := item2.id

	// Execute: Insert new item BETWEEN item1 and item2
	// Desired order: item1, newItem, item2
	newItemText := "New Item (should be in middle)"
	newItemLevel := 1
	newItemSumFactor := 1
	newItemDataSource, _ := data_source.FromString("sum")

	params := []UpdateReportParamsItem{
		{
			ItemId: &item1Id, // Position 0 → positionMap[item1Id] = 1
		},
		{
			ItemId:     nil, // Position 1 → NOT in positionMap (new item)
			Text:       &newItemText,
			Level:      &newItemLevel,
			SumFactor:  &newItemSumFactor,
			DataSource: &newItemDataSource,
		},
		{
			ItemId: &item2Id, // Position 2 → positionMap[item2Id] = 3
		},
	}

	err := section.SynchronizeItems(params)
	if err != nil {
		t.Fatalf("SynchronizeItems() error = %v, want nil", err)
	}

	// Verify
	if len(section.items) != 3 {
		t.Fatalf("Section should have 3 items, got %d", len(section.items))
	}

	// EXPECTED: item1, newItem, item2
	// ACTUAL: Let's see what we get...

	t.Logf("After SynchronizeItems:")
	for i, item := range section.items {
		t.Logf("  Position %d: %s (sequence: %d, id: %s)", i, item.text, item.sequence, item.id)
	}

	// Verify the EXPECTED order
	if section.items[0].id != item1Id {
		t.Errorf("Position 0 should be Item 1, got %s", section.items[0].text)
	}
	if section.items[1].text != "New Item (should be in middle)" {
		t.Errorf("Position 1 should be New Item, got %s", section.items[1].text)
	}
	if section.items[2].id != item2Id {
		t.Errorf("Position 2 should be Item 2, got %s", section.items[2].text)
	}

	// Verify sequences are sequential
	for i, item := range section.items {
		expectedSeq := i + 1
		if item.sequence != expectedSeq {
			t.Errorf("Item at position %d should have sequence %d, got %d", i, expectedSeq, item.sequence)
		}
	}
}

// TestSection_SynchronizeItems_AddMultipleNewItems tests adding multiple new items
// at different positions
func TestSection_SynchronizeItems_AddMultipleNewItems(t *testing.T) {
	// Setup: Section with 3 existing items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	item3 := createTestItem("Item 3", 3, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2, item3})

	item1Id := item1.id
	item2Id := item2.id
	item3Id := item3.id

	// Execute: Insert new items at various positions
	// Desired order: newItemA, item1, newItemB, item2, item3, newItemC
	newItemAText := "New Item A (beginning)"
	newItemBText := "New Item B (middle)"
	newItemCText := "New Item C (end)"
	level := 1
	sumFactor := 1
	ds, _ := data_source.FromString("sum")

	params := []UpdateReportParamsItem{
		// Position 0: newItemA (new)
		{
			ItemId:     nil,
			Text:       &newItemAText,
			Level:      &level,
			SumFactor:  &sumFactor,
			DataSource: &ds,
		},
		// Position 1: item1 (existing) → positionMap[item1Id] = 2
		{
			ItemId: &item1Id,
		},
		// Position 2: newItemB (new)
		{
			ItemId:     nil,
			Text:       &newItemBText,
			Level:      &level,
			SumFactor:  &sumFactor,
			DataSource: &ds,
		},
		// Position 3: item2 (existing) → positionMap[item2Id] = 4
		{
			ItemId: &item2Id,
		},
		// Position 4: item3 (existing) → positionMap[item3Id] = 5
		{
			ItemId: &item3Id,
		},
		// Position 5: newItemC (new)
		{
			ItemId:     nil,
			Text:       &newItemCText,
			Level:      &level,
			SumFactor:  &sumFactor,
			DataSource: &ds,
		},
	}

	err := section.SynchronizeItems(params)
	if err != nil {
		t.Fatalf("SynchronizeItems() error = %v, want nil", err)
	}

	// Verify
	if len(section.items) != 6 {
		t.Fatalf("Section should have 6 items, got %d", len(section.items))
	}

	t.Logf("After SynchronizeItems with multiple new items:")
	for i, item := range section.items {
		t.Logf("  Position %d: %s (sequence: %d)", i, item.text, item.sequence)
	}

	// Expected order: newItemA, item1, newItemB, item2, item3, newItemC
	expectedTexts := []string{
		"New Item A (beginning)",
		"Item 1",
		"New Item B (middle)",
		"Item 2",
		"Item 3",
		"New Item C (end)",
	}

	for i, expectedText := range expectedTexts {
		if section.items[i].text != expectedText {
			t.Errorf("Position %d: expected '%s', got '%s'", i, expectedText, section.items[i].text)
		}
		if section.items[i].sequence != i+1 {
			t.Errorf("Position %d: expected sequence %d, got %d", i, i+1, section.items[i].sequence)
		}
	}
}

// TestSection_SynchronizeItems_NewItemAtBeginning tests adding a new item at the beginning
func TestSection_SynchronizeItems_NewItemAtBeginning(t *testing.T) {
	// Setup: Section with 2 existing items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2})

	item1Id := item1.id
	item2Id := item2.id

	// Execute: Add new item at the BEGINNING
	newItemText := "New Item (at beginning)"
	level := 1
	sumFactor := 1
	ds, _ := data_source.FromString("sum")

	params := []UpdateReportParamsItem{
		// Position 0: newItem (new)
		{
			ItemId:     nil,
			Text:       &newItemText,
			Level:      &level,
			SumFactor:  &sumFactor,
			DataSource: &ds,
		},
		// Position 1: item1 → positionMap[item1Id] = 2
		{
			ItemId: &item1Id,
		},
		// Position 2: item2 → positionMap[item2Id] = 3
		{
			ItemId: &item2Id,
		},
	}

	err := section.SynchronizeItems(params)
	if err != nil {
		t.Fatalf("SynchronizeItems() error = %v, want nil", err)
	}

	// Verify
	if len(section.items) != 3 {
		t.Fatalf("Section should have 3 items, got %d", len(section.items))
	}

	t.Logf("After adding new item at beginning:")
	for i, item := range section.items {
		t.Logf("  Position %d: %s (sequence: %d)", i, item.text, item.sequence)
	}

	// Expected order: newItem, item1, item2
	if section.items[0].text != "New Item (at beginning)" {
		t.Errorf("Position 0 should be New Item, got %s", section.items[0].text)
	}
	if section.items[1].id != item1Id {
		t.Errorf("Position 1 should be Item 1, got %s", section.items[1].text)
	}
	if section.items[2].id != item2Id {
		t.Errorf("Position 2 should be Item 2, got %s", section.items[2].text)
	}
}

// TestSection_SynchronizeItems_NewItemAtEnd tests adding a new item at the end
func TestSection_SynchronizeItems_NewItemAtEnd(t *testing.T) {
	// Setup: Section with 2 existing items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2})

	item1Id := item1.id
	item2Id := item2.id

	// Execute: Add new item at the END
	newItemText := "New Item (at end)"
	level := 1
	sumFactor := 1
	ds, _ := data_source.FromString("sum")

	params := []UpdateReportParamsItem{
		// Position 0: item1 → positionMap[item1Id] = 1
		{
			ItemId: &item1Id,
		},
		// Position 1: item2 → positionMap[item2Id] = 2
		{
			ItemId: &item2Id,
		},
		// Position 2: newItem (new)
		{
			ItemId:     nil,
			Text:       &newItemText,
			Level:      &level,
			SumFactor:  &sumFactor,
			DataSource: &ds,
		},
	}

	err := section.SynchronizeItems(params)
	if err != nil {
		t.Fatalf("SynchronizeItems() error = %v, want nil", err)
	}

	// Verify
	if len(section.items) != 3 {
		t.Fatalf("Section should have 3 items, got %d", len(section.items))
	}

	t.Logf("After adding new item at end:")
	for i, item := range section.items {
		t.Logf("  Position %d: %s (sequence: %d)", i, item.text, item.sequence)
	}

	// Expected order: item1, item2, newItem
	if section.items[0].id != item1Id {
		t.Errorf("Position 0 should be Item 1, got %s", section.items[0].text)
	}
	if section.items[1].id != item2Id {
		t.Errorf("Position 1 should be Item 2, got %s", section.items[1].text)
	}
	if section.items[2].text != "New Item (at end)" {
		t.Errorf("Position 2 should be New Item, got %s", section.items[2].text)
	}
}

// TestSection_reorderItemsBySequence_DebugNewItemSequences tests the internal sequence
// assignment logic for new items to understand the bug
func TestSection_reorderItemsBySequence_DebugNewItemSequences(t *testing.T) {
	// Setup: Section with 2 existing items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2})

	item1Id := item1.id
	item2Id := item2.id

	// We'll manually simulate what SynchronizeItems does
	// Step 1: Build desiredById map and newItems list
	newItemText := "New Item"
	level := 1
	sumFactor := 1
	ds, _ := data_source.FromString("sum")

	params := []UpdateReportParamsItem{
		{ItemId: &item1Id}, // Position 0
		{
			ItemId:     nil, // Position 1 (new item)
			Text:       &newItemText,
			Level:      &level,
			SumFactor:  &sumFactor,
			DataSource: &ds,
		},
		{ItemId: &item2Id}, // Position 2
	}

	// Simulate SynchronizeItems logic
	desiredById := make(map[uuid.UUID]*UpdateReportParamsItem)
	var newItems []*UpdateReportParamsItem

	for i := range params {
		item := &params[i]
		if item.ItemId == nil {
			newItems = append(newItems, item)
		} else {
			desiredById[*item.ItemId] = item
		}
	}

	t.Logf("desiredById map contains %d items", len(desiredById))
	t.Logf("newItems list contains %d items", len(newItems))

	// Step 2: Process existing items (keptItems)
	var keptItems []*Item
	for _, currentItem := range section.items {
		if desiredItem, exists := desiredById[currentItem.id]; exists {
			if err := currentItem.applyUpdates(*desiredItem); err != nil {
				t.Fatal(err)
			}
			keptItems = append(keptItems, currentItem)
		}
	}

	t.Logf("After processing existing items, keptItems has %d items", len(keptItems))

	// Step 3: Create new items
	// This is where the sequence is assigned: len(keptItems) + i + 1
	for i, newItemData := range newItems {
		sequence := len(keptItems) + i + 1
		t.Logf("Creating new item with sequence: %d (len(keptItems)=%d, i=%d)", sequence, len(keptItems), i)

		newItem, err := section.createItemFromData(*newItemData, sequence)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("New item created with text='%s', sequence=%d", newItem.text, newItem.sequence)
		keptItems = append(keptItems, newItem)
	}

	// Step 4: Replace section's items
	section.items = keptItems

	t.Logf("\nBefore reordering:")
	for i, item := range section.items {
		t.Logf("  Position %d: %s (sequence: %d, id: %s)", i, item.text, item.sequence, item.id)
	}

	// Step 5: Build positionMap for reordering
	positionMap := make(map[uuid.UUID]int)
	for i, desired := range params {
		if desired.ItemId != nil {
			positionMap[*desired.ItemId] = i + 1
			t.Logf("positionMap[%s (%s)] = %d", desired.ItemId,
				func() string {
					if *desired.ItemId == item1Id {
						return "Item 1"
					} else if *desired.ItemId == item2Id {
						return "Item 2"
					}
					return "Unknown"
				}(), i+1)
		}
	}

	t.Logf("\nNow testing the sorting logic:")
	// Let's manually check what the sort comparison would return
	// Between Item 1 and New Item:
	item1Pos, item1Has := positionMap[item1.id]
	newItemPos, newItemHas := positionMap[section.items[2].id] // New item is at index 2 in keptItems
	t.Logf("Comparing Item 1 vs New Item:")
	t.Logf("  Item 1: positionMap has=%v, position=%d", item1Has, item1Pos)
	t.Logf("  New Item: positionMap has=%v, position=%d, current sequence=%d", newItemHas, newItemPos, section.items[2].sequence)

	// Between New Item and Item 2:
	item2Pos, item2Has := positionMap[item2.id]
	t.Logf("Comparing New Item vs Item 2:")
	t.Logf("  New Item: positionMap has=%v, current sequence=%d", newItemHas, section.items[2].sequence)
	t.Logf("  Item 2: positionMap has=%v, position=%d", item2Has, item2Pos)
	t.Logf("  Comparison (newItem.sequence < item2Pos): %d < %d = %v", section.items[2].sequence, item2Pos, section.items[2].sequence < item2Pos)

	// Step 6: Call reorderItemsBySequence
	if err := section.reorderItemsBySequence(params); err != nil {
		t.Fatal(err)
	}

	t.Logf("\nAfter reorderItemsBySequence:")
	for i, item := range section.items {
		t.Logf("  Position %d: %s (sequence: %d)", i, item.text, item.sequence)
	}

	// Verify expected order
	expectedTexts := []string{"Item 1", "New Item", "Item 2"}
	for i, expectedText := range expectedTexts {
		if section.items[i].text != expectedText {
			t.Errorf("Position %d: expected '%s', got '%s'", i, expectedText, section.items[i].text)
		}
	}
}
