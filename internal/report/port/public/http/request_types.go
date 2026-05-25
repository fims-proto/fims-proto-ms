package http

import (
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/report/app/command"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"

	"github.com/google/uuid"
)

type UpdateReportRequest struct {
	Title   *string               `json:"title,omitempty"`
	Columns []UpdateColumnRequest `json:"columns,omitempty"`
	Rows    []UpdateRowRequest    `json:"rows"`
}

type UpdateColumnRequest struct {
	Id        *string `json:"id,omitempty"`
	Label     string  `json:"label"`
	ValueType string  `json:"valueType"`
}

type UpdateRowRequest struct {
	Id          *string                 `json:"id,omitempty"`
	RowCode     string                  `json:"rowCode"`
	Text        string                  `json:"text"`
	LineNo      *int                    `json:"lineNo,omitempty"`
	ShowLineNo  bool                    `json:"showLineNo"`
	SumFactor   int                     `json:"sumFactor"`
	CanEdit     *bool                   `json:"canEdit,omitempty"`
	CanMove     *bool                   `json:"canMove,omitempty"`
	CanAddChild *bool                   `json:"canAddChild,omitempty"`
	Expression  UpdateExpressionRequest `json:"expression"`
	Rows        []UpdateRowRequest      `json:"rows,omitempty"`
}

type UpdateExpressionRequest struct {
	Kind           string                                `json:"kind"`
	LedgerAccounts []UpdateLedgerAccountReferenceRequest `json:"ledgerAccounts,omitempty"`
	CashFlowItems  []UpdateCashFlowItemReferenceRequest  `json:"cashFlowItems,omitempty"`
	RowReferences  []UpdateRowReferenceRequest           `json:"rowReferences,omitempty"`
}

type UpdateLedgerAccountReferenceRequest struct {
	RawAccountNumber string    `json:"rawAccountNumber"`
	AccountId        uuid.UUID `json:"accountId,omitempty"`
	SumFactor        int       `json:"sumFactor"`
	Measure          string    `json:"measure"`
}

type UpdateCashFlowItemReferenceRequest struct {
	Code      string    `json:"code"`
	ItemId    uuid.UUID `json:"itemId,omitempty"`
	SumFactor int       `json:"sumFactor"`
}

type UpdateRowReferenceRequest struct {
	RowCode   string `json:"rowCode"`
	SumFactor int    `json:"sumFactor"`
}

func (r UpdateReportRequest) mapToCommand(reportId uuid.UUID, sobId uuid.UUID) (command.UpdateReportCmd, error) {
	columns := make([]command.UpdateReportCmdColumn, 0, len(r.Columns))
	for _, columnReq := range r.Columns {
		columnId, err := parseOptionalUUID(columnReq.Id, "columnId")
		if err != nil {
			return command.UpdateReportCmd{}, err
		}
		columns = append(columns, command.UpdateReportCmdColumn{
			ColumnId:  columnId,
			Label:     columnReq.Label,
			ValueType: columnReq.ValueType,
		})
	}

	rows, err := convertRows(r.Rows)
	if err != nil {
		return command.UpdateReportCmd{}, err
	}

	return command.UpdateReportCmd{
		ReportId: reportId,
		SobId:    sobId,
		Title:    r.Title,
		Columns:  columns,
		Rows:     rows,
	}, nil
}

func convertRows(reqs []UpdateRowRequest) ([]command.UpdateReportCmdRow, error) {
	rows := make([]command.UpdateReportCmdRow, 0, len(reqs))
	for _, rowReq := range reqs {
		rowId, err := parseOptionalUUID(rowReq.Id, "rowId")
		if err != nil {
			return nil, err
		}
		childRows, err := convertRows(rowReq.Rows)
		if err != nil {
			return nil, err
		}
		rows = append(rows, command.UpdateReportCmdRow{
			RowId:       rowId,
			RowCode:     rowReq.RowCode,
			Text:        rowReq.Text,
			LineNo:      rowReq.LineNo,
			ShowLineNo:  rowReq.ShowLineNo,
			SumFactor:   rowReq.SumFactor,
			CanEdit:     rowReq.CanEdit,
			CanMove:     rowReq.CanMove,
			CanAddChild: rowReq.CanAddChild,
			Expression: command.UpdateReportCmdExpression{
				Kind:           rowReq.Expression.Kind,
				LedgerAccounts: convertLedgerAccountRefs(rowReq.Expression.LedgerAccounts),
				CashFlowItems:  convertCashFlowItemRefs(rowReq.Expression.CashFlowItems),
				RowReferences:  convertRowRefs(rowReq.Expression.RowReferences),
			},
			Rows: childRows,
		})
	}
	return rows, nil
}

func convertLedgerAccountRefs(reqs []UpdateLedgerAccountReferenceRequest) []report.LedgerAccountReference {
	refs := make([]report.LedgerAccountReference, 0, len(reqs))
	for _, req := range reqs {
		refs = append(refs, report.LedgerAccountReference{
			RawAccountNumber: req.RawAccountNumber,
			AccountId:        req.AccountId,
			SumFactor:        req.SumFactor,
			Measure:          req.Measure,
		})
	}
	return refs
}

func convertCashFlowItemRefs(reqs []UpdateCashFlowItemReferenceRequest) []report.CashFlowItemReference {
	refs := make([]report.CashFlowItemReference, 0, len(reqs))
	for _, req := range reqs {
		refs = append(refs, report.CashFlowItemReference{
			Code:      req.Code,
			ItemId:    req.ItemId,
			SumFactor: req.SumFactor,
		})
	}
	return refs
}

func convertRowRefs(reqs []UpdateRowReferenceRequest) []report.RowReference {
	refs := make([]report.RowReference, 0, len(reqs))
	for _, req := range reqs {
		refs = append(refs, report.RowReference{
			RowCode:   req.RowCode,
			SumFactor: req.SumFactor,
		})
	}
	return refs
}

func parseOptionalUUID(value *string, field string) (uuid.UUID, error) {
	if value == nil || *value == "" {
		return uuid.Nil, nil
	}
	parsed, err := uuid.Parse(*value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s: %s", field, *value)
	}
	return parsed, nil
}
