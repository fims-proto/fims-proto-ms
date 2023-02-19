package http

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type AccountResponse struct {
	Id                uuid.UUID `json:"id"`
	SobId             uuid.UUID `json:"sobId"`
	SuperiorAccountId uuid.UUID `json:"superiorAccountId"`
	Title             string    `json:"title"`
	AccountNumber     string    `json:"accountNumber"`
	NumberHierarchy   []int     `json:"numberHierarchy"`
	Level             int       `json:"level"`
	AccountType       string    `json:"accountType"`
	BalanceDirection  string    `json:"balanceDirection"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type PeriodResponse struct {
	Id           uuid.UUID `json:"id"`
	SobId        uuid.UUID `json:"sobId"`
	FiscalYear   int       `json:"fiscalYear"`
	PeriodNumber int       `json:"periodNumber"`
	OpeningTime  time.Time `json:"openingTime"`
	EndingTime   time.Time `json:"endingTime"`
	IsClosed     bool      `json:"isClosed"`
	IsCurrent    bool      `json:"isCurrent"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type LedgerResponse struct {
	Id             uuid.UUID       `json:"id"`
	SobId          uuid.UUID       `json:"sobId"`
	AccountId      uuid.UUID       `json:"accountId"`
	PeriodId       uuid.UUID       `json:"periodId"`
	OpeningBalance decimal.Decimal `json:"openingBalance"`
	EndingBalance  decimal.Decimal `json:"endingBalance"`
	PeriodDebit    decimal.Decimal `json:"periodDebit"`
	PeriodCredit   decimal.Decimal `json:"periodCredit"`
	Account        AccountResponse `json:"account"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

// mapper

func accountDTOToVO(dto query.Account) AccountResponse {
	return AccountResponse{
		Id:                dto.Id,
		SobId:             dto.SobId,
		SuperiorAccountId: dto.SuperiorAccountId,
		Title:             dto.Title,
		AccountNumber:     dto.AccountNumber,
		NumberHierarchy:   dto.NumberHierarchy,
		Level:             dto.Level,
		AccountType:       dto.AccountType,
		BalanceDirection:  dto.BalanceDirection,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
	}
}

func periodDTOToVO(dto query.Period) PeriodResponse {
	return PeriodResponse{
		Id:           dto.Id,
		SobId:        dto.SobId,
		FiscalYear:   dto.FiscalYear,
		PeriodNumber: dto.PeriodNumber,
		OpeningTime:  dto.OpeningTime,
		EndingTime:   dto.EndingTime,
		IsClosed:     dto.IsClosed,
		IsCurrent:    dto.IsCurrent,
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
	}
}

func ledgerDTOToVO(dto query.Ledger) LedgerResponse {
	return LedgerResponse{
		Id:             dto.Id,
		SobId:          dto.SobId,
		AccountId:      dto.AccountId,
		PeriodId:       dto.PeriodId,
		OpeningBalance: dto.OpeningBalance,
		EndingBalance:  dto.EndingBalance,
		PeriodDebit:    dto.PeriodDebit,
		PeriodCredit:   dto.PeriodCredit,
		Account:        accountDTOToVO(dto.Account),
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
	}
}
