package http

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type ReportResponse struct {
	Id          uuid.UUID         `json:"id,omitempty"`
	SobId       uuid.UUID         `json:"sobId,omitempty"`
	Period      *PeriodResponse   `json:"period,omitempty"`
	Title       string            `json:"title,omitempty"`
	Template    bool              `json:"template"`
	Class       string            `json:"class"`
	AmountTypes []string          `json:"amountTypes"`
	Sections    []SectionResponse `json:"sections"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PeriodResponse struct {
	FiscalYear   int `json:"fiscalYear"`
	PeriodNumber int `json:"periodNumber"`
}

type SectionResponse struct {
	Id       uuid.UUID         `json:"id,omitempty"`
	Title    string            `json:"title"`
	Amounts  []decimal.Decimal `json:"amounts,omitempty"`
	Sections []SectionResponse `json:"sections,omitempty"`
	Items    []ItemResponse    `json:"items,omitempty"`
}

type ItemResponse struct {
	Id               uuid.UUID         `json:"id,omitempty"`
	Text             string            `json:"text"`
	Level            int               `json:"level"`
	SumFactor        int               `json:"sumFactor"`
	DisplaySumFactor bool              `json:"displaySumFactor,omitempty"`
	DataSource       string            `json:"dataSource"`
	Formulas         []FormulaResponse `json:"formulas,omitempty"`
	Amounts          []decimal.Decimal `json:"amounts,omitempty"`
	IsBreakdownItem  bool              `json:"isBreakdownItem,omitempty"`
	IsDeletable      bool              `json:"isDeletable,omitempty"`
	IsTextModifiable bool              `json:"isTextModifiable,omitempty"`
	IsDraggable      bool              `json:"isDraggable,omitempty"`
	IsAbleToAddChild bool              `json:"isAbleToAddChild,omitempty"`
	IsAbleToAddLeaf  bool              `json:"isAbleToAddLeaf,omitempty"`
}

type FormulaResponse struct {
	Id        uuid.UUID         `json:"id"`
	Account   AccountResponse   `json:"account"`
	SumFactor int               `json:"sumFactor"`
	Rule      string            `json:"rule"`
	Amounts   []decimal.Decimal `json:"amounts,omitempty"`
}

type AccountResponse struct {
	Id                uuid.UUID  `json:"id,omitempty"`
	SobId             uuid.UUID  `json:"sobId,omitempty"`
	SuperiorAccountId *uuid.UUID `json:"superiorAccountId,omitempty"`
	Title             string     `json:"title"`
	AccountNumber     string     `json:"accountNumber"`
	Level             int        `json:"level"`
	IsLeaf            bool       `json:"isLeaf"`
	Class             int        `json:"class"`
	Group             int        `json:"group"`
	BalanceDirection  string     `json:"balanceDirection"`
}

// mappers

func reportDTOToVO(dto query.Report) ReportResponse {
	return ReportResponse{
		Id:          dto.Id,
		SobId:       dto.SobId,
		Period:      periodDTOToVO(dto.Period),
		Title:       dto.Title,
		Template:    dto.Template,
		Class:       dto.Class,
		AmountTypes: dto.AmountTypes,
		Sections:    converter.DTOsToVOs(dto.Sections, sectionDTOToVO),
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
}

func sectionDTOToVO(dto query.Section) SectionResponse {
	return SectionResponse{
		Id:       dto.Id,
		Title:    dto.Title,
		Amounts:  dto.Amounts,
		Sections: converter.DTOsToVOs(dto.Sections, sectionDTOToVO),
		Items:    converter.DTOsToVOs(dto.Items, itemDTOToVO),
	}
}

func itemDTOToVO(dto query.Item) ItemResponse {
	return ItemResponse{
		Id:               dto.Id,
		Text:             dto.Text,
		Level:            dto.Level,
		SumFactor:        dto.SumFactor,
		DisplaySumFactor: dto.DisplaySumFactor,
		DataSource:       dto.DataSource,
		Formulas:         converter.DTOsToVOs(dto.Formulas, formulaDTOToVO),
		Amounts:          dto.Amounts,
		IsBreakdownItem:  dto.IsBreakdownItem,
		IsDeletable:      dto.IsDeletable,
		IsTextModifiable: dto.IsTextModifiable,
		IsDraggable:      dto.IsDraggable,
		IsAbleToAddChild: dto.IsAbleToAddChild,
		IsAbleToAddLeaf:  dto.IsAbleToAddLeaf,
	}
}

func formulaDTOToVO(dto query.Formula) FormulaResponse {
	return FormulaResponse{
		Id:        dto.Id,
		Account:   accountDTOToVO(dto.Account),
		SumFactor: dto.SumFactor,
		Rule:      dto.Rule,
		Amounts:   dto.Amounts,
	}
}

func accountDTOToVO(dto query.Account) AccountResponse {
	return AccountResponse{
		Id:                dto.Id,
		SobId:             dto.SobId,
		SuperiorAccountId: dto.SuperiorAccountId,
		Title:             dto.Title,
		AccountNumber:     dto.AccountNumber,
		Level:             dto.Level,
		IsLeaf:            dto.IsLeaf,
		Class:             dto.Class,
		Group:             dto.Group,
		BalanceDirection:  dto.BalanceDirection,
	}
}

func periodDTOToVO(dto *query.Period) *PeriodResponse {
	if dto == nil {
		return nil
	}
	return &PeriodResponse{
		FiscalYear:   dto.FiscalYear,
		PeriodNumber: dto.PeriodNumber,
	}
}
