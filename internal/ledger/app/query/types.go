package query

import (
	"time"

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
	AccountId      uuid.UUID
	OpeningBalance decimal.Decimal
	EndingBalance  decimal.Decimal
	Debit          decimal.Decimal
	Credit         decimal.Decimal
	CreatedAt      time.Time
	UpdatedAt      time.Time
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
