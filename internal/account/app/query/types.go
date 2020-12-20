package query

import "github/fims-proto/fims-proto-ms/internal/account/domain/account"

type Account struct {
	number         string
	title          string
	superiorNumber string
	accountType    account.Type
}
