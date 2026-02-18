package account

import (
	"fmt"
	"unicode/utf8"
)

func (a *Account) UpdateTitle(title string) error {
	if title == "" {
		return fmt.Errorf("empty title")
	}

	if utf8.RuneCountInString(title) > 50 {
		return fmt.Errorf("account title exceeds max length (50)")
	}

	a.title = title
	return nil
}
