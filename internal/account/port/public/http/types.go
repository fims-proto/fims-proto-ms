package http

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type AccountConfigurationResponse struct {
	SobId             string    `json:"sobId"`
	AccountId         string    `json:"accountId"`
	SuperiorAccountId string    `json:"superiorAccountId"`
	Title             string    `json:"title"`
	AccountNumber     string    `json:"accountNumber"`
	NumberHierarchy   []int     `json:"numberHierarchy"`
	Level             int       `json:"level"`
	AccountType       string    `json:"accountType"`
	BalanceDirection  string    `json:"balanceDirection"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func accountConfigurationDTOToVO(q query.AccountConfiguration) AccountConfigurationResponse {
	return AccountConfigurationResponse{
		SobId:             q.SobId.String(),
		AccountId:         q.AccountId.String(),
		SuperiorAccountId: q.SuperiorAccountId.String(),
		Title:             q.Title,
		AccountNumber:     q.AccountNumber,
		NumberHierarchy:   q.NumberHierarchy,
		Level:             q.Level,
		AccountType:       q.AccountType,
		BalanceDirection:  q.BalanceDirection,
		CreatedAt:         q.CreatedAt,
		UpdatedAt:         q.UpdatedAt,
	}
}
