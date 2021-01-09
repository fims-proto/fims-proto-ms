package query

import (
	"github/fims-proto/fims-proto-ms/internal/account/domain/account_type"
)

type Account struct {
	number         string
	title          string
	superiorNumber string
	accountType    accounttype.Type
}
