package http

import (
	"github/fims-proto/fims-proto-ms/internal/sob/app/command"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
)

type slugErr interface {
	Slug() string
}

type Error struct {
	Message string
	Slug    string
}

type CreateSobRequest struct {
	Name                string
	Description         string
	BaseCurrency        string
	StartingPeriodYear  int
	StartingPeriodMonth int
	AccountsCodeLength  []int
}

type UpdateSobRequest struct {
	Name               string
	AccountsCodeLength []int
}

type SobResponse struct {
	Id                  string
	Name                string
	Description         string
	BaseCurrency        string
	StartingPeriodYear  int
	StartingPeriodMonth int
	AccountsCodeLength  []int
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
