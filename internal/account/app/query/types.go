package query

import (
	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
)

type Account struct {
	Id                uuid.UUID
	SobId             uuid.UUID
	SuperiorAccountId uuid.UUID
	AccountNumber     string
	NumberHierarchy   []int
	Title             string
	AccountType       commonAccount.Type
	BalanceDirection  commonAccount.Direction
	SuperiorAccount   *Account
}
