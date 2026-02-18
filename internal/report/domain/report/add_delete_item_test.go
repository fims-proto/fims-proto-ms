package report

import (
	"testing"

	"github.com/google/uuid"
)

func TestReport_AddItemToSection(t *testing.T) {
	// Create a test report with a section
	section1 := createTestSection("Section 1", 1, nil, nil)
	if section1 == nil {
		t.Fatal("createTestSection returned nil")
	}
	sectionId := section1.id
	report := createTestReport([]*Section{section1})
	if report == nil {
		t.Fatal("createTestReport returned nil")
	}

	// Create a new item to add
	newItem := createTestItem("New Item", 1, true)

	// Test adding item to existing section (default: insert at beginning)
	err := report.AddItemToSection(sectionId, newItem, 0)
	if err != nil {
		t.Errorf("AddItemToSection() error = %v, want nil", err)
	}

	// Verify item was added - need to find the section again since it's modified in place
	foundSection, _ := report.findSectionById(sectionId)
	if len(foundSection.items) != 1 {
		t.Errorf("Section should have 1 item, got %d", len(foundSection.items))
	}

	// Verify sequence was set correctly
	if foundSection.items[0].sequence != 1 {
		t.Errorf("Item sequence should be 1, got %d", foundSection.items[0].sequence)
	}

	// Test adding to non-existent section
	err = report.AddItemToSection(uuid.New(), newItem, 0)
	if err == nil {
		t.Error("AddItemToSection() should return error for non-existent section")
	}

	// Test adding nil item
	err = report.AddItemToSection(sectionId, nil, 0)
	if err == nil {
		t.Error("AddItemToSection() should return error for nil item")
	}
}

func TestReport_DeleteItemFromSection(t *testing.T) {
	// Create items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	item3 := createTestItem("Item 3", 3, false)
	item2Id := item2.id
	item3Id := item3.id

	// Create section with items
	section1 := createTestSection("Section 1", 1, nil, []*Item{item1, item2, item3})
	sectionId := section1.id
	report := createTestReport([]*Section{section1})

	// Test deleting an item
	err := report.DeleteItemFromSection(sectionId, item2Id)
	if err != nil {
		t.Errorf("DeleteItemFromSection() error = %v, want nil", err)
	}

	// Verify item was deleted
	foundSection, _ := report.findSectionById(sectionId)
	if len(foundSection.items) != 2 {
		t.Errorf("Section should have 2 items, got %d", len(foundSection.items))
	}

	// Verify sequences were renumbered (1, 3 -> 1, 2)
	if foundSection.items[0].sequence != 1 {
		t.Errorf("First item sequence should be 1, got %d", foundSection.items[0].sequence)
	}
	if foundSection.items[1].sequence != 2 {
		t.Errorf("Second item sequence should be 2, got %d", foundSection.items[1].sequence)
	}

	// Test deleting non-editable item
	err = report.DeleteItemFromSection(sectionId, item3Id)
	if err == nil {
		t.Error("DeleteItemFromSection() should return error for non-editable item")
	}

	// Test deleting non-existent item in valid section
	err = report.DeleteItemFromSection(sectionId, uuid.New())
	if err == nil {
		t.Error("DeleteItemFromSection() should return error for non-existent item")
	}

	// Test deleting item with non-existent section ID
	err = report.DeleteItemFromSection(uuid.New(), item2Id)
	if err == nil {
		t.Error("DeleteItemFromSection() should return error for non-existent section")
	}
}

func TestSection_AddItem(t *testing.T) {
	section := createTestSection("Test Section", 1, nil, nil)

	// Add first item (insert at beginning, default behavior)
	item1 := createTestItem("Item 1", 999, true)
	err := section.AddItem(item1, 0)
	if err != nil {
		t.Errorf("AddItem() error = %v, want nil", err)
	}

	if len(section.items) != 1 {
		t.Errorf("Section should have 1 item, got %d", len(section.items))
	}
	if section.items[0].sequence != 1 {
		t.Errorf("Item sequence should be 1, got %d", section.items[0].sequence)
	}

	// Add second item (insert at beginning again)
	item2 := createTestItem("Item 2", 999, true)
	err = section.AddItem(item2, 0)
	if err != nil {
		t.Errorf("AddItem() error = %v, want nil", err)
	}

	if len(section.items) != 2 {
		t.Errorf("Section should have 2 items, got %d", len(section.items))
	}
	// item2 should be first now
	if section.items[0].id != item2.id || section.items[0].sequence != 1 {
		t.Errorf("First item should be item2 with sequence 1")
	}
	if section.items[1].id != item1.id || section.items[1].sequence != 2 {
		t.Errorf("Second item should be item1 with sequence 2")
	}
}

func TestSection_DeleteItem(t *testing.T) {
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	item3 := createTestItem("Item 3", 3, true)

	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2, item3})

	// Delete middle item
	err := section.DeleteItem(item2.id)
	if err != nil {
		t.Errorf("DeleteItem() error = %v, want nil", err)
	}

	if len(section.items) != 2 {
		t.Errorf("Section should have 2 items, got %d", len(section.items))
	}

	// Verify renumbering
	if section.items[0].id != item1.id || section.items[0].sequence != 1 {
		t.Errorf("First item should be item1 with sequence 1")
	}
	if section.items[1].id != item3.id || section.items[1].sequence != 2 {
		t.Errorf("Second item should be item3 with sequence 2 (was 3)")
	}
}

func TestSection_renumberItems(t *testing.T) {
	item1 := createTestItem("Item 1", 5, true)
	item2 := createTestItem("Item 2", 10, true)
	item3 := createTestItem("Item 3", 15, true)

	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2, item3})

	section.renumberItems()

	// Verify all items have sequential numbers
	expectedSequences := []int{1, 2, 3}
	for i, item := range section.items {
		if item.sequence != expectedSequences[i] {
			t.Errorf("Item %d should have sequence %d, got %d", i, expectedSequences[i], item.sequence)
		}
	}
}

func TestSection_findSectionByIdRecursive(t *testing.T) {
	// Create nested sections
	childSection := createTestSection("Child Section", 1, nil, nil)
	parentSection := createTestSection("Parent Section", 1, []*Section{childSection}, nil)

	// Find child section
	found := parentSection.findSectionByIdRecursive(childSection.id)
	if found == nil {
		t.Error("Should find child section")
		return
	}
	if found.id != childSection.id {
		t.Error("Found wrong section")
	}

	// Find parent section
	found = parentSection.findSectionByIdRecursive(parentSection.id)
	if found == nil {
		t.Error("Should find parent section")
		return
	}
	if found.id != parentSection.id {
		t.Error("Found wrong section")
	}

	// Find non-existent section
	found = parentSection.findSectionByIdRecursive(uuid.New())
	if found != nil {
		t.Error("Should not find non-existent section")
	}
}

// Helper functions for creating test objects

func createTestReport(sections []*Section) *Report {
	report, err := New(
		uuid.New(),
		uuid.New(),
		uuid.New(),
		"Test Report",
		false,
		"balance_sheet",
		[]string{"period_ending_balance"},
		sections,
	)
	if err != nil {
		panic(err)
	}
	return report
}

func createTestSection(title string, sequence int, childSections []*Section, items []*Item) *Section {
	if childSections == nil {
		childSections = []*Section{}
	}
	if items == nil {
		items = []*Item{}
	}
	section, err := NewSection(
		uuid.New(),
		title,
		sequence,
		"", // Empty string for None section type
		nil,
		childSections,
		items,
	)
	if err != nil {
		panic(err)
	}
	return section
}

func createTestItem(text string, sequence int, isEditable bool) *Item {
	item, err := NewItem(
		uuid.New(),
		text,
		1,
		sequence,
		"",
		1,
		false,
		"sum",
		nil,
		nil,
		isEditable,
		false,
		false,
	)
	if err != nil {
		panic(err)
	}
	return item
}

// Tests for insert position control

func TestSection_AddItem_InsertAtBeginning(t *testing.T) {
	// Create section with existing items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	item3 := createTestItem("Item 3", 3, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2, item3})

	// Insert new item at beginning (insertAfterSequence = 0)
	newItem := createTestItem("New Item", 999, true)
	err := section.AddItem(newItem, 0)
	if err != nil {
		t.Errorf("AddItem() error = %v, want nil", err)
	}

	// Verify order and sequences
	if len(section.items) != 4 {
		t.Errorf("Section should have 4 items, got %d", len(section.items))
	}

	expectedOrder := []*Item{newItem, item1, item2, item3}
	expectedSequences := []int{1, 2, 3, 4}

	for i, expectedItem := range expectedOrder {
		if section.items[i].id != expectedItem.id {
			t.Errorf("Item at position %d should be %s, got %s", i, expectedItem.text, section.items[i].text)
		}
		if section.items[i].sequence != expectedSequences[i] {
			t.Errorf("Item at position %d should have sequence %d, got %d", i, expectedSequences[i], section.items[i].sequence)
		}
	}
}

func TestSection_AddItem_InsertInMiddle(t *testing.T) {
	// Create section with existing items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	item3 := createTestItem("Item 3", 3, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2, item3})

	// Insert new item after sequence 2 (in the middle)
	newItem := createTestItem("New Item", 999, true)
	err := section.AddItem(newItem, 2)
	if err != nil {
		t.Errorf("AddItem() error = %v, want nil", err)
	}

	// Verify order and sequences
	// Expected: item1(1), item2(2), newItem(3), item3(4)
	if len(section.items) != 4 {
		t.Errorf("Section should have 4 items, got %d", len(section.items))
	}

	expectedOrder := []*Item{item1, item2, newItem, item3}
	expectedSequences := []int{1, 2, 3, 4}

	for i, expectedItem := range expectedOrder {
		if section.items[i].id != expectedItem.id {
			t.Errorf("Item at position %d should be %s, got %s", i, expectedItem.text, section.items[i].text)
		}
		if section.items[i].sequence != expectedSequences[i] {
			t.Errorf("Item at position %d should have sequence %d, got %d", i, expectedSequences[i], section.items[i].sequence)
		}
	}
}

func TestSection_AddItem_InsertAtEnd(t *testing.T) {
	// Create section with existing items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2})

	// Insert new item at end (insertAfterSequence >= max)
	newItem := createTestItem("New Item", 999, true)
	err := section.AddItem(newItem, 999)
	if err != nil {
		t.Errorf("AddItem() error = %v, want nil", err)
	}

	// Verify order and sequences
	if len(section.items) != 3 {
		t.Errorf("Section should have 3 items, got %d", len(section.items))
	}

	expectedOrder := []*Item{item1, item2, newItem}
	expectedSequences := []int{1, 2, 3}

	for i, expectedItem := range expectedOrder {
		if section.items[i].id != expectedItem.id {
			t.Errorf("Item at position %d should be %s, got %s", i, expectedItem.text, section.items[i].text)
		}
		if section.items[i].sequence != expectedSequences[i] {
			t.Errorf("Item at position %d should have sequence %d, got %d", i, expectedSequences[i], section.items[i].sequence)
		}
	}
}

func TestSection_AddItem_InsertIntoEmptySection(t *testing.T) {
	section := createTestSection("Test Section", 1, nil, nil)

	// Insert item with various insertAfterSequence values - should all result in sequence=1
	testCases := []int{0, -1, 1, 999}

	for _, insertAfter := range testCases {
		// Reset section
		section.items = []*Item{}

		newItem := createTestItem("New Item", 999, true)
		err := section.AddItem(newItem, insertAfter)
		if err != nil {
			t.Errorf("AddItem(insertAfter=%d) error = %v, want nil", insertAfter, err)
		}

		if len(section.items) != 1 {
			t.Errorf("Section should have 1 item, got %d", len(section.items))
		}

		if section.items[0].sequence != 1 {
			t.Errorf("Item sequence should be 1 for insertAfter=%d, got %d", insertAfter, section.items[0].sequence)
		}
	}
}

func TestSection_AddItem_NegativeInsertAfterSequence(t *testing.T) {
	// Create section with existing items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2})

	// Insert with negative insertAfterSequence should insert at beginning
	newItem := createTestItem("New Item", 999, true)
	err := section.AddItem(newItem, -5)
	if err != nil {
		t.Errorf("AddItem() error = %v, want nil", err)
	}

	// Verify newItem is at the beginning
	if section.items[0].id != newItem.id {
		t.Error("New item should be at the beginning")
	}
	if section.items[0].sequence != 1 {
		t.Errorf("New item sequence should be 1, got %d", section.items[0].sequence)
	}
}

func TestSection_AddItem_InsertAfterFirst(t *testing.T) {
	// Create section with existing items
	item1 := createTestItem("Item 1", 1, true)
	item2 := createTestItem("Item 2", 2, true)
	item3 := createTestItem("Item 3", 3, true)
	section := createTestSection("Test Section", 1, nil, []*Item{item1, item2, item3})

	// Insert after first item (insertAfterSequence = 1)
	newItem := createTestItem("New Item", 999, true)
	err := section.AddItem(newItem, 1)
	if err != nil {
		t.Errorf("AddItem() error = %v, want nil", err)
	}

	// Verify order: item1(1), newItem(2), item2(3), item3(4)
	expectedOrder := []*Item{item1, newItem, item2, item3}
	expectedSequences := []int{1, 2, 3, 4}

	for i, expectedItem := range expectedOrder {
		if section.items[i].id != expectedItem.id {
			t.Errorf("Item at position %d should be %s, got %s", i, expectedItem.text, section.items[i].text)
		}
		if section.items[i].sequence != expectedSequences[i] {
			t.Errorf("Item at position %d should have sequence %d, got %d", i, expectedSequences[i], section.items[i].sequence)
		}
	}
}
