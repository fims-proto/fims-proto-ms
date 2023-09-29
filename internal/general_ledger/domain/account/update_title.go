package account

import "fmt"

func (a *Account) UpdateTitle(title string) error {
	if title == "" {
		return fmt.Errorf("empty title")
	}

	a.title = title
	return nil
}
