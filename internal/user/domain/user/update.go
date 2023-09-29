package user

import (
	"encoding/json"
	"errors"
)

func (u *User) Update(traits json.RawMessage) error {
	if len(traits) == 0 {
		return errors.New("traits is empty")
	}

	u.traits = traits
	return nil
}
