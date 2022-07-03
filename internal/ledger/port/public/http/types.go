package http

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type PeriodResponse struct {
	Id               uuid.UUID `json:"id"`
	SobId            uuid.UUID `json:"sobId"`
	PreviousPeriodId uuid.UUID `json:"previousPeriodId"`
	FinancialYear    int       `json:"financialYear"`
	Number           int       `json:"number"`
	OpeningTime      time.Time `json:"openingTime"`
	EndingTime       time.Time `json:"endingTime"`
	IsClosed         bool      `json:"isClosed"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CreatePeriodRequest struct {
	PreviousPeriodId uuid.UUID `json:"previousPeriodId"`
	FinancialYear    int       `json:"financialYear"`
	Number           int       `json:"number"`
	OpeningTime      time.Time `json:"openingTime"`
	EndingTime       time.Time `json:"endingTime"`
}

type LedgerResponse struct {
	Id             uuid.UUID       `json:"id"`
	PeriodId       uuid.UUID       `json:"periodId"`
	Account        AccountResponse `json:"account"`
	OpeningBalance decimal.Decimal `json:"openingBalance"`
	EndingBalance  decimal.Decimal `json:"endingBalance"`
	Debit          decimal.Decimal `json:"debit"`
	Credit         decimal.Decimal `json:"credit"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type AccountResponse struct {
	Id                string `json:"id"`
	SuperiorAccountId string `json:"superiorAccountId"`
	AccountNumber     string `json:"accountNumber"`
	Title             string `json:"title"`
	AccountType       string `json:"accountType"`
	BalanceDirection  string `json:"balanceDirection"`
}

func mapFromPeriodQuery(p query.Period) PeriodResponse {
	return PeriodResponse{
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
		Id:       l.Id,
		PeriodId: l.PeriodId,
		Account: AccountResponse{
			Id:                l.Account.Id.String(),
			SuperiorAccountId: l.Account.SuperiorAccountId.String(),
			AccountNumber:     l.Account.AccountNumber,
			Title:             l.Account.Title,
			AccountType:       l.Account.AccountType.String(),
			BalanceDirection:  l.Account.BalanceDirection.String(),
		},
		OpeningBalance: l.OpeningBalance,
		EndingBalance:  l.EndingBalance,
		Debit:          l.Debit,
		Credit:         l.Credit,
		CreatedAt:      l.CreatedAt,
		UpdatedAt:      l.UpdatedAt,
	}
}

func (r CreatePeriodRequest) mapToCommand() command.CreatePeriodCmd {
	return command.CreatePeriodCmd{
		PreviousPeriodId: r.PreviousPeriodId,
		SobId:            uuid.Nil,
		FinancialYear:    r.FinancialYear,
		Number:           r.Number,
		OpeningTime:      r.OpeningTime,
		EndingTime:       r.EndingTime,
	}
}
