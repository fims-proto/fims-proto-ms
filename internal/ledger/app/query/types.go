package query

import (
	"time"

	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountingPeriod struct {
	Id               uuid.UUID
	SobId            uuid.UUID
	PreviousPeriodId uuid.UUID
	FinancialYear    int
	Number           int
	OpeningTime      time.Time
	EndingTime       time.Time
	IsClosed         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Ledger struct {
	Id             uuid.UUID
	PeriodId       uuid.UUID
	Account        Account
	OpeningBalance decimal.Decimal
	EndingBalance  decimal.Decimal
	Debit          decimal.Decimal
	Credit         decimal.Decimal
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Account struct {
	Id                uuid.UUID
	SuperiorAccountId uuid.UUID
	AccountNumber     string
	Title             string
	AccountType       commonAccount.Type
	BalanceDirection  commonAccount.Direction
}

type LedgerLog struct {
	Id              uuid.UUID
	PostingId       uuid.UUID
	AccountId       uuid.UUID
	VoucherId       uuid.UUID
	TransactionTime time.Time
	Debit           decimal.Decimal
	Credit          decimal.Decimal
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
