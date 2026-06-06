package sob

import (
	"fmt"
	"unicode/utf8"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (s *Sob) UpdateName(name string) error {
	if name == "" {
		return commonErrors.NewInvalidInputError(commonErrors.SlugSobEmptyName)
	}

	if utf8.RuneCountInString(name) > 50 {
		return commonErrors.NewInvalidInputError(commonErrors.SlugSobNameTooLong)
	}

	s.name = name
	return nil
}

func (s *Sob) UpdateDescription(description string) error {
	if utf8.RuneCountInString(description) > 500 {
		return commonErrors.NewInvalidInputError(commonErrors.SlugSobDescriptionTooLong)
	}

	s.description = description
	return nil
}

func (s *Sob) UpdateAccountsCodeLength(accountsCodeLength []int) error {
	if len(accountsCodeLength) < 2 || len(accountsCodeLength) > 10 {
		return commonErrors.NewInvalidInputError(commonErrors.SlugSobInvalidAccountLevel)
	}

	// account level can only be enlarged
	if len(accountsCodeLength) < len(s.accountsCodeLength) {
		return commonErrors.NewInvalidInputError(commonErrors.SlugSobCannotShortenLevel)
	}

	for i, accountCodeLength := range accountsCodeLength {
		if accountCodeLength < 1 || accountCodeLength > 6 {
			return fmt.Errorf("invalid account code length at level %d", i)
		}
		// account code length can only be enlarged
		if i < len(s.accountsCodeLength) && accountCodeLength < s.accountsCodeLength[i] {
			return fmt.Errorf("cannot reduce account code length at level %d", i)
		}
	}

	s.accountsCodeLength = accountsCodeLength
	return nil
}
