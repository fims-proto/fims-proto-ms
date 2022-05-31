package query

import (
	"time"

	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
)

type Account struct {
	Id                uuid.UUID
	SobId             uuid.UUID
	Title             string
	AccountNumber     string
	NumberHierarchy   []int
	SuperiorAccountId uuid.UUID
	AccountType       commonAccount.Type
	BalanceDirection  commonAccount.Direction
	Level             int
	SuperiorAccount   *Account
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
