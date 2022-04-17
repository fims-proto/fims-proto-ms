package http

import "github/fims-proto/fims-proto-ms/internal/account/app/query"

type slugErr interface {
	Slug() string
}

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type AccountResponse struct {
	Id                string `json:"id"`
	SobId             string `json:"sobId"`
	SuperiorAccountId string `json:"superiorAccountId"`
	AccountNumber     string `json:"accountNumber"`
	Title             string `json:"title"`
	Level             int    `json:"level"`
	AccountType       string `json:"accountType"`
	BalanceDirection  string `json:"balanceDirection"`
}

func mapFromAccountQuery(q query.Account) AccountResponse {
	return AccountResponse{
		Id:                q.Id.String(),
		SobId:             q.SobId.String(),
		SuperiorAccountId: q.SuperiorAccountId.String(),
		AccountNumber:     q.AccountNumber,
		Title:             q.Title,
		Level:             q.Level,
		AccountType:       q.AccountType.String(),
		BalanceDirection:  q.BalanceDirection.String(),
	}
}
