package query

import (
	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
)

type Account struct {
	Id                uuid.UUID
	SobId             uuid.UUID
	SuperiorAccountId uuid.UUID
	SuperiorNumbers   []int
	LevelNumber       int
	Title             string
	Level             uint8
	AccountType       string
	BalanceDirection  commonAccount.Direction
	SuperiorAccount   *Account
}
