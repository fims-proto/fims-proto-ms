package user

import (
	"encoding/json"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

func (u *User) Update(traits json.RawMessage) error {
	if len(traits) == 0 {
		return commonErrors.NewInvalidInputError(commonErrors.SlugUserEmptyTraits)
	}

	u.traits = traits
	return nil
}
