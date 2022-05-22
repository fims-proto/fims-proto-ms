package domain

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func (u *User) Update(traits json.RawMessage) error {
	if len(traits) == 0 {
		return errors.New("traits is empty")
	}

	u.traits = traits
	return nil
}
