package class

import (
	"cmp"
	"slices"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

func Validate(c Class, g Group) error {
	i, found := slices.BinarySearchFunc(Classes, pair{Class: c}, func(a pair, b pair) int {
		return cmp.Compare(a.Class, b.Class)
	})

	if !found {
		return errors.NewInvalidInputError(errors.SlugInvalidAccountClass, c.String())
	}

	if !slices.Contains(Classes[i].Groups, g) {
		return errors.NewInvalidInputError(errors.SlugInvalidAccountGroup, c.String(), g.String())
	}

	return nil
}
