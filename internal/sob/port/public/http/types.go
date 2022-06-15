package http

import (
	"github/fims-proto/fims-proto-ms/internal/sob/app/command"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type CreateSobRequest struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	BaseCurrency        string `json:"baseCurrency"`
	StartingPeriodYear  int    `json:"startingPeriodYear"`
	StartingPeriodMonth int    `json:"startingPeriodMonth"`
	AccountsCodeLength  []int  `json:"accountsCodeLength"`
}

type UpdateSobRequest struct {
	Name               string `json:"name"`
	AccountsCodeLength []int  `json:"accountsCodeLength"`
}

type SobResponse struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	BaseCurrency        string `json:"baseCurrency"`
	StartingPeriodYear  int    `json:"startingPeriodYear"`
	StartingPeriodMonth int    `json:"startingPeriodMonth"`
	AccountsCodeLength  []int  `json:"accountsCodeLength"`
}

func mapFromSobQuery(q query.Sob) SobResponse {
	return SobResponse{
		Id:                  q.Id.String(),
		Name:                q.Name,
		Description:         q.Description,
		BaseCurrency:        q.BaseCurrency,
		StartingPeriodYear:  q.StartingPeriodYear,
		StartingPeriodMonth: q.StartingPeriodMonth,
		AccountsCodeLength:  q.AccountsCodeLength,
	}
}

func (r CreateSobRequest) mapToCommand() command.CreateSobCmd {
	return command.CreateSobCmd{
		Name:                r.Name,
		Description:         r.Description,
		BaseCurrency:        r.BaseCurrency,
		StartingPeriodYear:  r.StartingPeriodYear,
		StartingPeriodMonth: r.StartingPeriodMonth,
		AccountsCodeLength:  r.AccountsCodeLength,
	}
}
