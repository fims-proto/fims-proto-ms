package http

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"
)

type Error struct {
	Message string `json:"message,omitempty"`
	Slug    string `json:"slug,omitempty"`
}

type AccountConfigurationResponse struct {
	SobId             uuid.UUID `json:"sobId,omitempty"`
	AccountId         uuid.UUID `json:"accountId,omitempty"`
	SuperiorAccountId uuid.UUID `json:"superiorAccountId,omitempty"`
	Title             string    `json:"title,omitempty"`
	AccountNumber     string    `json:"accountNumber,omitempty"`
	NumberHierarchy   []int     `json:"numberHierarchy,omitempty"`
	Level             int       `json:"level,omitempty"`
	AccountType       string    `json:"accountType,omitempty"`
	BalanceDirection  string    `json:"balanceDirection,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type PeriodResponse struct {
	SobId            uuid.UUID `json:"sobId,omitempty"`
	PeriodId         uuid.UUID `json:"periodId,omitempty"`
	PreviousPeriodId uuid.UUID `json:"previousPeriodId,omitempty"`
	FinancialYear    int       `json:"financialYear,omitempty"`
	Number           int       `json:"number,omitempty"`
	OpeningTime      time.Time `json:"openingTime"`
	EndingTime       time.Time `json:"endingTime"`
	IsClosed         bool      `json:"isClosed,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type AccountResponse struct {
	SobId          uuid.UUID                    `json:"sobId,omitempty"`
	AccountId      uuid.UUID                    `json:"accountId,omitempty"`
	PeriodId       uuid.UUID                    `json:"periodId"`
	OpeningBalance decimal.Decimal              `json:"openingBalance"`
	EndingBalance  decimal.Decimal              `json:"endingBalance"`
	PeriodDebit    decimal.Decimal              `json:"periodDebit"`
	PeriodCredit   decimal.Decimal              `json:"periodCredit"`
	Configuration  AccountConfigurationResponse `json:"configuration"`
	CreatedAt      time.Time                    `json:"createdAt"`
	UpdatedAt      time.Time                    `json:"updatedAt"`
}

// mapper

func accountConfigurationDTOToVO(dto query.AccountConfiguration) AccountConfigurationResponse {
	return AccountConfigurationResponse{
		SobId:             dto.SobId,
		AccountId:         dto.AccountId,
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
		SobId:            dto.SobId,
		PeriodId:         dto.PeriodId,
		PreviousPeriodId: dto.PreviousPeriodId,
		FinancialYear:    dto.FinancialYear,
		Number:           dto.Number,
		OpeningTime:      dto.OpeningTime,
		EndingTime:       dto.EndingTime,
		IsClosed:         dto.IsClosed,
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
	}
}

func accountDTOToVO(dto query.Account) AccountResponse {
	return AccountResponse{
		SobId:          dto.SobId,
		AccountId:      dto.AccountId,
		PeriodId:       dto.PeriodId,
		OpeningBalance: dto.OpeningBalance,
		EndingBalance:  dto.EndingBalance,
		PeriodDebit:    dto.PeriodDebit,
		PeriodCredit:   dto.PeriodCredit,
		Configuration:  accountConfigurationDTOToVO(dto.Configuration),
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
	}
}
