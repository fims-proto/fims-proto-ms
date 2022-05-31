package http

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"
)

type slugErr interface {
	Slug() string
}

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type AccountResponse struct {
	Id                string    `json:"id"`
	SobId             string    `json:"sobId"`
	AccountNumber     string    `json:"accountNumber"`
	Title             string    `json:"title"`
	NumberHierarchy   []int     `json:"numberHierarchy"`
	SuperiorAccountId string    `json:"superiorAccountId"`
	AccountType       string    `json:"accountType"`
	BalanceDirection  string    `json:"balanceDirection"`
	Level             int       `json:"level"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func mapFromAccountQuery(q query.Account) AccountResponse {
	return AccountResponse{
		Id:                q.Id.String(),
		SobId:             q.SobId.String(),
		AccountNumber:     q.AccountNumber,
		Title:             q.Title,
		NumberHierarchy:   q.NumberHierarchy,
		SuperiorAccountId: q.SuperiorAccountId.String(),
		AccountType:       q.AccountType.String(),
		BalanceDirection:  q.BalanceDirection.String(),
		Level:             q.Level,
		CreatedAt:         q.CreatedAt,
		UpdatedAt:         q.UpdatedAt,
	}
}
