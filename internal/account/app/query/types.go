package query

import (
	"time"

	"github.com/shopspring/decimal"

	"github.com/google/uuid"
)

type Account struct {
	Id                uuid.UUID
	SobId             uuid.UUID
	SuperiorAccountId uuid.UUID
	Title             string
	AccountNumber     string
	NumberHierarchy   []int
	Level             int
	AccountType       string
	BalanceDirection  string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Period struct {
	SobId            uuid.UUID
	PeriodId         uuid.UUID
	PreviousPeriodId uuid.UUID
	FiscalYear       int
	PeriodNumber     int
	OpeningTime      time.Time
	EndingTime       time.Time
	IsClosed         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Ledger struct {
	Id             uuid.UUID
	SobId          uuid.UUID
	AccountId      uuid.UUID
	PeriodId       uuid.UUID
	OpeningBalance decimal.Decimal
	EndingBalance  decimal.Decimal
	PeriodDebit    decimal.Decimal
	PeriodCredit   decimal.Decimal
	Account        Account
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
