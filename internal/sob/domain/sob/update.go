package sob

import (
	"errors"
	"fmt"
)

func (s *Sob) UpdateName(name string) error {
	if name == "" {
		return errors.New("empty sob name")
	}

	s.name = name
	return nil
}

func (s *Sob) UpdateAccountsCodeLength(accountsCodeLength []int) error {
	if len(accountsCodeLength) < 2 || len(accountsCodeLength) > 10 {
		return errors.New("invalid account level")
	}

	// account level can only be enlarged
	if len(accountsCodeLength) < len(s.accountsCodeLength) {
		return errors.New("cannot shorten account level")
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
