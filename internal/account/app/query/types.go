package query

import (
	"github.com/shopspring/decimal"
	"time"

	"github.com/google/uuid"
)

type AccountConfiguration struct {
	SobId             uuid.UUID
	AccountId         uuid.UUID
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
	FinancialYear    int
	Number           int
	OpeningTime      time.Time
	EndingTime       time.Time
	IsClosed         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Account struct {
	SobId          uuid.UUID
	AccountId      uuid.UUID
	OpeningBalance decimal.Decimal
	EndingBalance  decimal.Decimal
	PeriodDebit    decimal.Decimal
	PeriodCredit   decimal.Decimal
	Configuration  AccountConfiguration
	Period         Period
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
