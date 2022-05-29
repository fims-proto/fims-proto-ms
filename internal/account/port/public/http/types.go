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
	NumberHierarchy   []int  `json:"numberHierarchy"`
	Title             string `json:"title"`
	AccountType       string `json:"accountType"`
	BalanceDirection  string `json:"balanceDirection"`
}

func mapFromAccountQuery(q query.Account) AccountResponse {
	return AccountResponse{
		Id:                q.Id.String(),
		SobId:             q.SobId.String(),
		SuperiorAccountId: q.SuperiorAccountId.String(),
		AccountNumber:     q.AccountNumber,
		NumberHierarchy:   q.NumberHierarchy,
		Title:             q.Title,
		AccountType:       q.AccountType.String(),
		BalanceDirection:  q.BalanceDirection.String(),
	}
}
