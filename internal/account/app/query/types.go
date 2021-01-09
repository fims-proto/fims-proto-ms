package query

import (
	"github/fims-proto/fims-proto-ms/internal/account/domain/account_type"
)

type Account struct {
	Number         string
	Title          string
	SuperiorNumber string
	AccountType    accounttype.Type
}
