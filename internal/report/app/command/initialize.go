package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"

	"github.com/google/uuid"
)

type InitializeHandler struct {
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService
}

func NewInitializeHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) InitializeHandler {
	if repo == nil {
		panic("nil repository")
	}
	if generalLedgerService == nil {
		panic("nil general ledger service")
	}
	return InitializeHandler{repo: repo, generalLedgerService: generalLedgerService}
}

func (h InitializeHandler) Handle(ctx context.Context, cmd InitializeCmd) error {
	reports, err := h.loadReports(ctx, cmd.SobId)
	if err != nil {
		return err
	}
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.repo.CreateReports(txCtx, reports)
	})
}

func (h InitializeHandler) loadReports(ctx context.Context, sobId uuid.UUID) ([]*report.Report, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get working directory: %w", err)
	}

	files := []string{
		"report_balance_sheet.json",
		"report_income_statement.json",
		"report_cash_flow_statement.json",
	}

	var reports []*report.Report
	for _, file := range files {
		raw, err := os.ReadFile(filepath.Join(workDir, "dataload", "xqykjzz", file))
		if err != nil {
			return nil, fmt.Errorf("could not read %s: %w", file, err)
		}
		var cmd initializeReport
		if err = json.Unmarshal(raw, &cmd); err != nil {
			return nil, fmt.Errorf("could not unmarshal %s: %w", file, err)
		}
		r, err := h.convertReport(ctx, sobId, cmd)
		if err != nil {
			return nil, fmt.Errorf("could not convert %s: %w", file, err)
		}
		reports = append(reports, r)
	}
	return reports, nil
}

type initializeReport struct {
	Title   string             `json:"title"`
	Class   string             `json:"class"`
	Columns []initializeColumn `json:"columns"`
	Rows    []initializeRow    `json:"rows"`
}

type initializeColumn struct {
	Label     string `json:"label"`
	ValueType string `json:"valueType"`
}

type initializeRow struct {
	RowCode     string               `json:"rowCode"`
	Text        string               `json:"text"`
	LineNo      *int                 `json:"lineNo"`
	ShowLineNo  bool                 `json:"showLineNo"`
	SumFactor   int                  `json:"sumFactor"`
	CanEdit     *bool                `json:"canEdit"`
	CanMove     *bool                `json:"canMove"`
	CanAddChild *bool                `json:"canAddChild"`
	Expression  initializeExpression `json:"expression"`
	Rows        []initializeRow      `json:"rows"`
}

type initializeExpression struct {
	Kind           string                          `json:"kind"`
	LedgerAccounts []report.LedgerAccountReference `json:"ledgerAccounts"`
	CashFlowItems  []report.CashFlowItemReference  `json:"cashFlowItems"`
	RowReferences  []report.RowReference           `json:"rowReferences"`
}

func (h InitializeHandler) convertReport(ctx context.Context, sobId uuid.UUID, cmd initializeReport) (*report.Report, error) {
	columns := make([]*report.Column, 0, len(cmd.Columns))
	for i, columnCmd := range cmd.Columns {
		column, err := report.NewColumn(uuid.New(), columnCmd.Label, columnCmd.ValueType, i+1)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	if err := h.resolveExpressionReferences(ctx, sobId, cmd.Rows); err != nil {
		return nil, err
	}

	rows, err := convertRows(cmd.Rows)
	if err != nil {
		return nil, err
	}

	return report.New(uuid.New(), sobId, uuid.Nil, cmd.Title, true, cmd.Class, columns, rows)
}

func (h InitializeHandler) resolveExpressionReferences(ctx context.Context, sobId uuid.UUID, rows []initializeRow) error {
	accountNumbers := make(map[string]struct{})
	cashFlowCodes := make(map[string]struct{})
	collectRefsFromRows(rows, accountNumbers, cashFlowCodes)

	var rawAccountNumbers []string
	for rawNumber := range accountNumbers {
		rawAccountNumbers = append(rawAccountNumbers, rawNumber)
	}
	accountIds := map[string]uuid.UUID{}
	if len(rawAccountNumbers) > 0 {
		var err error
		accountIds, err = h.generalLedgerService.ReadAccountIdsByRawNumbers(ctx, sobId, rawAccountNumbers)
		if err != nil {
			return fmt.Errorf("failed to resolve account ids: %w", err)
		}
	}

	var codes []string
	for code := range cashFlowCodes {
		codes = append(codes, code)
	}
	cashFlowItemIds := map[string]uuid.UUID{}
	if len(codes) > 0 {
		var err error
		cashFlowItemIds, err = h.generalLedgerService.ReadCashFlowItemIdsByCodes(ctx, sobId, codes)
		if err != nil {
			return fmt.Errorf("failed to resolve cash flow item ids: %w", err)
		}
	}

	return applyResolvedRefs(rows, accountIds, cashFlowItemIds)
}

func collectRefsFromRows(rows []initializeRow, accountNumbers map[string]struct{}, cashFlowCodes map[string]struct{}) {
	for ri := range rows {
		for _, ref := range rows[ri].Expression.LedgerAccounts {
			accountNumbers[ref.RawAccountNumber] = struct{}{}
		}
		for _, ref := range rows[ri].Expression.CashFlowItems {
			cashFlowCodes[ref.Code] = struct{}{}
		}
		collectRefsFromRows(rows[ri].Rows, accountNumbers, cashFlowCodes)
	}
}

func applyResolvedRefs(rows []initializeRow, accountIds map[string]uuid.UUID, cashFlowItemIds map[string]uuid.UUID) error {
	for ri := range rows {
		for li := range rows[ri].Expression.LedgerAccounts {
			ref := &rows[ri].Expression.LedgerAccounts[li]
			id, ok := accountIds[ref.RawAccountNumber]
			if !ok {
				return fmt.Errorf("account %s not found", ref.RawAccountNumber)
			}
			ref.AccountId = id
		}
		for ci := range rows[ri].Expression.CashFlowItems {
			ref := &rows[ri].Expression.CashFlowItems[ci]
			id, ok := cashFlowItemIds[ref.Code]
			if !ok {
				return fmt.Errorf("cash flow item %s not found", ref.Code)
			}
			ref.ItemId = id
		}
		if err := applyResolvedRefs(rows[ri].Rows, accountIds, cashFlowItemIds); err != nil {
			return err
		}
	}
	return nil
}

func convertRows(cmds []initializeRow) ([]*report.Row, error) {
	rows := make([]*report.Row, 0, len(cmds))
	for i, cmd := range cmds {
		expr, err := report.NewExpression(uuid.New(), cmd.Expression.Kind, cmd.Expression.LedgerAccounts, cmd.Expression.CashFlowItems, cmd.Expression.RowReferences)
		if err != nil {
			return nil, err
		}
		canEdit := boolValue(cmd.CanEdit, true)
		canMove := boolValue(cmd.CanMove, true)
		canAddChild := boolValue(cmd.CanAddChild, false)
		childRows, err := convertRows(cmd.Rows)
		if err != nil {
			return nil, err
		}
		row, err := report.NewRow(uuid.New(), cmd.RowCode, cmd.Text, i+1, cmd.LineNo, cmd.ShowLineNo, cmd.SumFactor, canEdit, canMove, canAddChild, expr, childRows, nil)
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func boolValue(v *bool, defaultValue bool) bool {
	if v == nil {
		return defaultValue
	}
	return *v
}
