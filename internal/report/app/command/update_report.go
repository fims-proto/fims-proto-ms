package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"

	"github.com/google/uuid"
)

type UpdateReportHandler struct {
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService
}

func NewUpdateReportHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) UpdateReportHandler {
	if repo == nil {
		panic("nil repository")
	}
	if generalLedgerService == nil {
		panic("nil general ledger service")
	}
	return UpdateReportHandler{repo: repo, generalLedgerService: generalLedgerService}
}

func (h UpdateReportHandler) Handle(ctx context.Context, cmd UpdateReportCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		if err := h.resolveExpressionReferences(txCtx, cmd.SobId, cmd.Rows); err != nil {
			return err
		}
		return h.repo.UpdateReport(txCtx, cmd.ReportId, func(r *report.Report) (*report.Report, error) {
			params, err := commandToUpdateParams(cmd)
			if err != nil {
				return nil, err
			}
			if err = r.UpdateStructure(params); err != nil {
				return nil, err
			}
			return r, nil
		})
	})
}

func (h UpdateReportHandler) resolveExpressionReferences(ctx context.Context, sobId uuid.UUID, rows []UpdateReportCmdRow) error {
	accountNumbers := make(map[string]struct{})
	cashFlowCodes := make(map[string]struct{})
	collectCmdRefs(rows, accountNumbers, cashFlowCodes)

	accountIds := map[string]uuid.UUID{}
	if len(accountNumbers) > 0 {
		var rawAccountNumbers []string
		for rawNumber := range accountNumbers {
			rawAccountNumbers = append(rawAccountNumbers, rawNumber)
		}
		var err error
		accountIds, err = h.generalLedgerService.ReadAccountIdsByRawNumbers(ctx, sobId, rawAccountNumbers)
		if err != nil {
			return fmt.Errorf("failed to resolve account ids: %w", err)
		}
	}

	cashFlowItemIds := map[string]uuid.UUID{}
	if len(cashFlowCodes) > 0 {
		var codes []string
		for code := range cashFlowCodes {
			codes = append(codes, code)
		}
		var err error
		cashFlowItemIds, err = h.generalLedgerService.ReadCashFlowItemIdsByCodes(ctx, sobId, codes)
		if err != nil {
			return fmt.Errorf("failed to resolve cash flow item ids: %w", err)
		}
	}

	return applyCmdRefs(rows, accountIds, cashFlowItemIds)
}

func collectCmdRefs(rows []UpdateReportCmdRow, accountNumbers map[string]struct{}, cashFlowCodes map[string]struct{}) {
	for ri := range rows {
		for _, ref := range rows[ri].Expression.LedgerAccounts {
			if ref.RawAccountNumber != "" {
				accountNumbers[ref.RawAccountNumber] = struct{}{}
			}
		}
		for _, ref := range rows[ri].Expression.CashFlowItems {
			if ref.Code != "" {
				cashFlowCodes[ref.Code] = struct{}{}
			}
		}
		collectCmdRefs(rows[ri].Rows, accountNumbers, cashFlowCodes)
	}
}

func applyCmdRefs(rows []UpdateReportCmdRow, accountIds map[string]uuid.UUID, cashFlowItemIds map[string]uuid.UUID) error {
	for ri := range rows {
		for li := range rows[ri].Expression.LedgerAccounts {
			ref := &rows[ri].Expression.LedgerAccounts[li]
			if ref.RawAccountNumber == "" {
				continue
			}
			id, ok := accountIds[ref.RawAccountNumber]
			if !ok {
				return fmt.Errorf("account %s not found", ref.RawAccountNumber)
			}
			ref.AccountId = id
		}
		for ci := range rows[ri].Expression.CashFlowItems {
			ref := &rows[ri].Expression.CashFlowItems[ci]
			if ref.Code == "" {
				continue
			}
			id, ok := cashFlowItemIds[ref.Code]
			if !ok {
				return fmt.Errorf("cash flow item %s not found", ref.Code)
			}
			ref.ItemId = id
		}
		if err := applyCmdRefs(rows[ri].Rows, accountIds, cashFlowItemIds); err != nil {
			return err
		}
	}
	return nil
}

func commandToUpdateParams(cmd UpdateReportCmd) (report.UpdateParams, error) {
	columns := make([]report.UpdateColumnParams, 0, len(cmd.Columns))
	for _, c := range cmd.Columns {
		columns = append(columns, report.UpdateColumnParams{
			ColumnId:  c.ColumnId,
			Label:     c.Label,
			ValueType: c.ValueType,
		})
	}

	rows, err := commandRowsToParams(cmd.Rows)
	if err != nil {
		return report.UpdateParams{}, err
	}
	return report.UpdateParams{Title: cmd.Title, Columns: columns, Rows: rows}, nil
}

func commandRowsToParams(cmds []UpdateReportCmdRow) ([]report.UpdateRowParams, error) {
	rows := make([]report.UpdateRowParams, 0, len(cmds))
	for _, rowCmd := range cmds {
		expr, err := report.NewExpression(uuid.New(), rowCmd.Expression.Kind, rowCmd.Expression.LedgerAccounts, rowCmd.Expression.CashFlowItems, rowCmd.Expression.RowReferences)
		if err != nil {
			return nil, err
		}
		childRows, err := commandRowsToParams(rowCmd.Rows)
		if err != nil {
			return nil, err
		}
		rows = append(rows, report.UpdateRowParams{
			RowId:       rowCmd.RowId,
			RowCode:     rowCmd.RowCode,
			Text:        rowCmd.Text,
			LineNo:      rowCmd.LineNo,
			ShowLineNo:  rowCmd.ShowLineNo,
			SumFactor:   rowCmd.SumFactor,
			CanEdit:     rowCmd.CanEdit,
			CanMove:     rowCmd.CanMove,
			CanAddChild: rowCmd.CanAddChild,
			Expression:  expr,
			Rows:        childRows,
		})
	}
	return rows, nil
}
