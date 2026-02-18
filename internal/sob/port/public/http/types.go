package http

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/sob/app/command"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"

	"github.com/google/uuid"
)

type Error struct {
	Message string `json:"message,omitempty"`
	Slug    string `json:"slug,omitempty"`
}

type CreateSobRequest struct {
	Name                string `json:"name,omitempty"`
	Description         string `json:"description,omitempty"`
	BaseCurrency        string `json:"baseCurrency,omitempty"`
	StartingPeriodYear  int    `json:"startingPeriodYear,omitempty"`
	StartingPeriodMonth int    `json:"startingPeriodMonth,omitempty"`
	AccountsCodeLength  []int  `json:"accountsCodeLength,omitempty"`
}

type UpdateSobRequest struct {
	Name               string  `json:"name,omitempty"`
	Description        *string `json:"description,omitempty"`
	AccountsCodeLength []int   `json:"accountsCodeLength,omitempty"`
}

type SobResponse struct {
	Id                  uuid.UUID `json:"id,omitempty"`
	Name                string    `json:"name,omitempty"`
	Description         string    `json:"description,omitempty"`
	BaseCurrency        string    `json:"baseCurrency,omitempty"`
	StartingPeriodYear  int       `json:"startingPeriodYear,omitempty"`
	StartingPeriodMonth int       `json:"startingPeriodMonth,omitempty"`
	AccountsCodeLength  []int     `json:"accountsCodeLength,omitempty"`
	CreatedAt           time.Time `json:"createdAt,omitempty"`
	UpdatedAt           time.Time `json:"updatedAt,omitempty"`
}

func sobDTOToVO(q query.Sob) SobResponse {
	return SobResponse(q)
}

func (r CreateSobRequest) mapToCommand() command.CreateSobCmd {
	return command.CreateSobCmd{
		SobId:               uuid.New(),
		Name:                r.Name,
		Description:         r.Description,
		BaseCurrency:        r.BaseCurrency,
		StartingPeriodYear:  r.StartingPeriodYear,
		StartingPeriodMonth: r.StartingPeriodMonth,
		AccountsCodeLength:  r.AccountsCodeLength,
	}
}
