package http

import (
	"github/fims-proto/fims-proto-ms/internal/ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type slugErr interface {
	Slug() string
}

type Error struct {
	Message string
	Slug    string
}

type AccountingPeriodResponse struct {
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

type CreateAccountingPeriodRequest struct {
	PreviousPeriodId uuid.UUID
	FinancialYear    int
	Number           int
	OpeningTime      time.Time
	EndingTime       time.Time
}

type LedgerResponse struct {
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

func mapFromPeriodQuery(p query.AccountingPeriod) AccountingPeriodResponse {
	return AccountingPeriodResponse{
		Id:               p.Id,
		SobId:            p.SobId,
		PreviousPeriodId: p.PreviousPeriodId,
		FinancialYear:    p.FinancialYear,
		Number:           p.Number,
		OpeningTime:      p.OpeningTime,
		EndingTime:       p.EndingTime,
		IsClosed:         p.IsClosed,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
}

func mapFromLedgerQuery(l query.Ledger) LedgerResponse {
	return LedgerResponse{
		Id:             l.Id,
		PeriodId:       l.PeriodId,
		AccountId:      l.AccountId,
		OpeningBalance: l.OpeningBalance,
		EndingBalance:  l.EndingBalance,
		Debit:          l.Debit,
		Credit:         l.Credit,
		CreatedAt:      l.CreatedAt,
		UpdatedAt:      l.UpdatedAt,
	}
}

func (r CreateAccountingPeriodRequest) mapToCommand() command.CreatePeriodCmd {
	return command.CreatePeriodCmd{
		PreviousPeriodId: r.PreviousPeriodId,
		SobId:            uuid.Nil,
		FinancialYear:    r.FinancialYear,
		Number:           r.Number,
		OpeningTime:      r.OpeningTime,
		EndingTime:       r.EndingTime,
	}
}
