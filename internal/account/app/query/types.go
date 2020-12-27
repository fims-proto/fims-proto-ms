package query

import (
	"github/fims-proto/fims-proto-ms/internal/account/domain/type"
)

type Account struct {
	number         string
	title          string
	superiorNumber string
	accountType    _type.Type
}
