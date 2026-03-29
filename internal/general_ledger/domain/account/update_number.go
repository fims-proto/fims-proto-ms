package account

import "fmt"

func (a *Account) UpdateNumber(levelNumber int) error {
	// Get current hierarchy and replace the last element with the new level number
	currentHierarchy, err := HierarchyFromRaw(a.RawAccountNumber())
	if err != nil {
		return fmt.Errorf("failed to extract hierarchy from raw account number: %w", err)
	}
	newHierarchy := currentHierarchy[:len(currentHierarchy)-1]
	newHierarchy = append(newHierarchy, levelNumber)

	// Create new raw account number from updated hierarchy
	newRawNumber, err := PadRawAccountNumber(newHierarchy)
	if err != nil {
		return fmt.Errorf("failed to pad raw account number: %w", err)
	}

	a.rawAccountNumber = newRawNumber
	return nil
}
