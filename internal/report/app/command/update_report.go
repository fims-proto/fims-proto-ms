package command

import (
	"context"
	"fmt"
	"maps"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/formula_rule"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/item_type"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"

	"github.com/google/uuid"
)

type UpdateReportCmd struct {
	ReportId    uuid.UUID
	SobId       uuid.UUID
	Title       *string
	AmountTypes []amount_type.AmountType
	Sections    []UpdateSectionData
}

type UpdateSectionData struct {
	SectionId uuid.UUID
	Title     *string
	Items     []UpdateItemData // Complete desired item list
}

type UpdateItemData struct {
	// Nil ID means create new item
	ItemId *uuid.UUID

	// Position
	Sequence int

	// Item content (for new items or updates)
	Text             *string
	Level            *int
	SumFactor        *int
	DisplaySumFactor *bool
	ItemType         *item_type.ItemType
	DataSource       *data_source.DataSource
	Formulas         []FormulaData
	IsBreakdownItem  *bool
	IsAbleToAddChild *bool
}

type FormulaData struct {
	SumFactor     int
	AccountNumber string
	AccountId     uuid.UUID // Resolved from AccountNumber
	Rule          formula_rule.FormulaRule
}

type UpdateReportHandler struct {
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService
}

func NewUpdateReportHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) UpdateReportHandler {
	if repo == nil {
		panic("nil repo")
	}

	if generalLedgerService == nil {
		panic("nil general ledger service")
	}

	return UpdateReportHandler{
		repo:                 repo,
		generalLedgerService: generalLedgerService,
	}
}

func (h UpdateReportHandler) Handle(ctx context.Context, cmd UpdateReportCmd) (map[string]string, error) {
	createdItemIds := make(map[string]string)

	err := h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.repo.UpdateReport(txCtx, cmd.ReportId, func(r *report.Report) (*report.Report, error) {
			// Resolve account numbers to account IDs for formulas
			if err := h.resolveAccountIds(txCtx, cmd.SobId, &cmd); err != nil {
				return nil, err
			}

			// Convert command to domain params
			params := h.cmdToParams(cmd)

			// Apply comprehensive update via domain method
			ids, err := r.UpdateReportStructure(params)
			if err != nil {
				return nil, err
			}

			// Capture created item IDs for response
			maps.Copy(createdItemIds, ids)

			return r, nil
		})
	})

	return createdItemIds, err
}

func (h UpdateReportHandler) resolveAccountIds(ctx context.Context, sobId uuid.UUID, cmd *UpdateReportCmd) error {
	// Collect all account numbers that need resolution
	accountNumbersSet := make(map[string]bool)
	for i := range cmd.Sections {
		for j := range cmd.Sections[i].Items {
			item := &cmd.Sections[i].Items[j]
			for k := range item.Formulas {
				accountNumbersSet[item.Formulas[k].AccountNumber] = true
			}
		}
	}

	// Convert set to slice
	var accountNumbers []string
	for accountNumber := range accountNumbersSet {
		accountNumbers = append(accountNumbers, accountNumber)
	}

	// Batch resolve all account numbers
	if len(accountNumbers) > 0 {
		accountIds, err := h.generalLedgerService.ReadAccountIdsByNumbers(ctx, sobId, accountNumbers)
		if err != nil {
			return fmt.Errorf("failed to read account ids by numbers: %w", err)
		}

		// Update all formulas with resolved account IDs
		for i := range cmd.Sections {
			for j := range cmd.Sections[i].Items {
				item := &cmd.Sections[i].Items[j]
				for k := range item.Formulas {
					formula := &item.Formulas[k]
					accountId, ok := accountIds[formula.AccountNumber]
					if !ok {
						return errors.NewSlugError("account-notFound", map[string]interface{}{
							"accountNumber": formula.AccountNumber,
						})
					}
					formula.AccountId = accountId
				}
			}
		}
	}

	return nil
}

func (h UpdateReportHandler) cmdToParams(cmd UpdateReportCmd) report.UpdateReportParams {
	var sections []report.UpdateSectionParams
	for _, sectionCmd := range cmd.Sections {
		var items []report.UpdateItemParams
		for _, itemCmd := range sectionCmd.Items {
			var formulas []report.UpdateFormulaParams
			for _, formulaCmd := range itemCmd.Formulas {
				formulas = append(formulas, report.UpdateFormulaParams{
					SumFactor: formulaCmd.SumFactor,
					AccountId: formulaCmd.AccountId,
					Rule:      formulaCmd.Rule,
				})
			}

			items = append(items, report.UpdateItemParams{
				ItemId:           itemCmd.ItemId,
				Sequence:         itemCmd.Sequence,
				Text:             itemCmd.Text,
				Level:            itemCmd.Level,
				SumFactor:        itemCmd.SumFactor,
				DisplaySumFactor: itemCmd.DisplaySumFactor,
				ItemType:         itemCmd.ItemType,
				DataSource:       itemCmd.DataSource,
				Formulas:         formulas,
				IsBreakdownItem:  itemCmd.IsBreakdownItem,
				IsAbleToAddChild: itemCmd.IsAbleToAddChild,
			})
		}

		sections = append(sections, report.UpdateSectionParams{
			SectionId: sectionCmd.SectionId,
			Title:     sectionCmd.Title,
			Items:     items,
		})
	}

	return report.UpdateReportParams{
		Title:       cmd.Title,
		AmountTypes: cmd.AmountTypes,
		Sections:    sections,
	}
}
