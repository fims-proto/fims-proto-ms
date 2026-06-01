package http

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type ReportResponse struct {
	Id        uuid.UUID        `json:"id,omitempty"`
	SobId     uuid.UUID        `json:"sobId,omitempty"`
	Period    *PeriodResponse  `json:"period,omitempty"`
	Title     string           `json:"title,omitempty"`
	Template  bool             `json:"template"`
	Class     string           `json:"class"`
	Columns   []ColumnResponse `json:"columns"`
	Rows      []RowResponse    `json:"rows"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
}

type PeriodResponse struct {
	FiscalYear   int `json:"fiscalYear"`
	PeriodNumber int `json:"periodNumber"`
}

type ColumnResponse struct {
	Id        uuid.UUID `json:"id,omitempty"`
	Label     string    `json:"label"`
	ValueType string    `json:"valueType"`
	Sequence  int       `json:"sequence"`
}

type RowResponse struct {
	Id               uuid.UUID          `json:"id,omitempty"`
	RowCode          string             `json:"rowCode"`
	Text             string             `json:"text"`
	Sequence         int                `json:"sequence"`
	LineNo           *int               `json:"lineNo,omitempty"`
	ShowLineNo       bool               `json:"showLineNo"`
	SumFactor        int                `json:"sumFactor"`
	DisplaySumFactor bool               `json:"displaySumFactor"`
	Indent           int                `json:"indent"`
	CanEdit          bool               `json:"canEdit"`
	CanMove          bool               `json:"canMove"`
	CanAddChild      bool               `json:"canAddChild"`
	Expression       ExpressionResponse `json:"expression"`
	Rows             []RowResponse      `json:"rows,omitempty"`
	Amounts          []decimal.Decimal  `json:"amounts,omitempty"`
}

type ExpressionResponse struct {
	Id             uuid.UUID                        `json:"id,omitempty"`
	Kind           string                           `json:"kind"`
	LedgerAccounts []LedgerAccountReferenceResponse `json:"ledgerAccounts,omitempty"`
	CashFlowItems  []CashFlowItemReferenceResponse  `json:"cashFlowItems,omitempty"`
	RowReferences  []RowReferenceResponse           `json:"rowReferences,omitempty"`
}

type LedgerAccountReferenceResponse struct {
	RawAccountNumber string    `json:"rawAccountNumber"`
	AccountId        uuid.UUID `json:"accountId,omitempty"`
	SumFactor        int       `json:"sumFactor"`
	Measure          string    `json:"measure"`
}

type CashFlowItemReferenceResponse struct {
	Code      string    `json:"code"`
	ItemId    uuid.UUID `json:"itemId,omitempty"`
	SumFactor int       `json:"sumFactor"`
}

type RowReferenceResponse struct {
	RowCode   string `json:"rowCode"`
	SumFactor int    `json:"sumFactor"`
}

func reportDTOToVO(dto query.Report) ReportResponse {
	return ReportResponse{
		Id:        dto.Id,
		SobId:     dto.SobId,
		Period:    periodDTOToVO(dto.Period),
		Title:     dto.Title,
		Template:  dto.Template,
		Class:     dto.Class,
		Columns:   converter.DTOsToVOs(dto.Columns, columnDTOToVO),
		Rows:      converter.DTOsToVOs(dto.Rows, rowDTOToVO),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func periodDTOToVO(dto *query.Period) *PeriodResponse {
	if dto == nil {
		return nil
	}
	return &PeriodResponse{FiscalYear: dto.FiscalYear, PeriodNumber: dto.PeriodNumber}
}

func columnDTOToVO(dto query.Column) ColumnResponse {
	return ColumnResponse{Id: dto.Id, Label: dto.Label, ValueType: dto.ValueType, Sequence: dto.Sequence}
}

func rowDTOToVO(dto query.Row) RowResponse {
	return RowResponse{
		Id:               dto.Id,
		RowCode:          dto.RowCode,
		Text:             dto.Text,
		Sequence:         dto.Sequence,
		LineNo:           dto.LineNo,
		ShowLineNo:       dto.ShowLineNo,
		SumFactor:        dto.SumFactor,
		DisplaySumFactor: dto.DisplaySumFactor,
		Indent:           dto.Indent,
		CanEdit:          dto.CanEdit,
		CanMove:          dto.CanMove,
		CanAddChild:      dto.CanAddChild,
		Expression:       expressionDTOToVO(dto.Expression),
		Rows:             converter.DTOsToVOs(dto.Rows, rowDTOToVO),
		Amounts:          dto.Amounts,
	}
}

func expressionDTOToVO(dto query.Expression) ExpressionResponse {
	return ExpressionResponse{
		Id:             dto.Id,
		Kind:           dto.Kind,
		LedgerAccounts: converter.DTOsToVOs(dto.LedgerAccounts, ledgerAccountRefDTOToVO),
		CashFlowItems:  converter.DTOsToVOs(dto.CashFlowItems, cashFlowItemRefDTOToVO),
		RowReferences:  converter.DTOsToVOs(dto.RowReferences, rowRefDTOToVO),
	}
}

func ledgerAccountRefDTOToVO(dto query.LedgerAccountReference) LedgerAccountReferenceResponse {
	return LedgerAccountReferenceResponse{
		RawAccountNumber: dto.RawAccountNumber,
		AccountId:        dto.AccountId,
		SumFactor:        dto.SumFactor,
		Measure:          dto.Measure,
	}
}

func cashFlowItemRefDTOToVO(dto query.CashFlowItemReference) CashFlowItemReferenceResponse {
	return CashFlowItemReferenceResponse{Code: dto.Code, ItemId: dto.ItemId, SumFactor: dto.SumFactor}
}

func rowRefDTOToVO(dto query.RowReference) RowReferenceResponse {
	return RowReferenceResponse{RowCode: dto.RowCode, SumFactor: dto.SumFactor}
}
