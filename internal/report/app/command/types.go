package command

import (
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"

	"github.com/google/uuid"
)

type InitializeCmd struct {
	SobId uuid.UUID
}

type GenerateReportCmd struct {
	SobId        uuid.UUID
	Class        string
	FiscalYear   int
	PeriodNumber int
}

type GenerateForPeriodCmd struct {
	SobId    uuid.UUID
	PeriodId uuid.UUID
}

type RecalculateReportCmd struct {
	ReportId uuid.UUID
}

type RegenerateReportCmd struct {
	ReportId uuid.UUID
}

type UpdateReportCmd struct {
	ReportId uuid.UUID
	SobId    uuid.UUID
	Title    *string
	Rows     []UpdateReportCmdRow
}

type UpdateReportCmdRow struct {
	RowId            uuid.UUID
	RowCode          string
	Text             string
	LineNo           *int
	ShowLineNo       bool
	SumFactor        int
	DisplaySumFactor bool
	Indent           int
	CanEdit          *bool
	CanMove          *bool
	CanAddChild      *bool
	Expression       UpdateReportCmdExpression
	Rows             []UpdateReportCmdRow
}

type UpdateReportCmdExpression struct {
	Kind           string
	LedgerAccounts []report.LedgerAccountReference
	CashFlowItems  []report.CashFlowItemReference
	RowReferences  []report.RowReference
}
