package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/amount_type"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/formula_rule"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"

	"github.com/google/uuid"
)

type UpdateReportCmd struct {
	ReportId    uuid.UUID
	SobId       uuid.UUID
	Title       *string
	AmountTypes []amount_type.AmountType
	Sections    []UpdateReportCmdSection
}

type UpdateReportCmdSection struct {
	SectionId uuid.UUID
	Title     *string
	Items     []UpdateReportCmdItem // Complete desired item list
	Sections  []UpdateReportCmdSection
}

type UpdateReportCmdItem struct {
	// Nil ID means create new item
	ItemId *uuid.UUID

	// Item content (for new items or updates)
	Text             *string
	Level            *int
	SumFactor        *int
	DisplaySumFactor *bool
	DataSource       *data_source.DataSource
	Formulas         []UpdateReportCmdFormula
	IsBreakdownItem  *bool
	IsAbleToAddChild *bool
}

type UpdateReportCmdFormula struct {
	FormulaId     *uuid.UUID
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

func (h UpdateReportHandler) Handle(ctx context.Context, cmd UpdateReportCmd) error {
	err := h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.repo.UpdateReport(txCtx, cmd.ReportId, func(r *report.Report) (*report.Report, error) {
			// Resolve account numbers to account IDs for formulas
			if err := h.resolveAccountIds(txCtx, cmd.SobId, &cmd); err != nil {
				return nil, err
			}

			// Convert command to domain params
			params := h.cmdToParams(cmd)

			// Apply comprehensive update via domain method
			err := r.UpdateReportStructure(params)
			if err != nil {
				return nil, err
			}

			return r, nil
		})
	})

	return err
}

func (h UpdateReportHandler) resolveAccountIds(ctx context.Context, sobId uuid.UUID, cmd *UpdateReportCmd) error {
	// Collect all account numbers that need resolution (recursively)
	accountNumbers := h.collectAccountNumbers(cmd.Sections)

	// Batch resolve all account numbers
	if len(accountNumbers) > 0 {
		accountIds, err := h.generalLedgerService.ReadAccountIdsByNumbers(ctx, sobId, accountNumbers)
		if err != nil {
			return fmt.Errorf("failed to read account ids by numbers: %w", err)
		}

		// Update all formulas with resolved account IDs (recursively)
		if err := h.updateFormulaAccountIds(cmd.Sections, accountIds); err != nil {
			return err
		}
	}

	return nil
}

// collectAccountNumbers recursively collects all unique account numbers from sections
func (h UpdateReportHandler) collectAccountNumbers(sections []UpdateReportCmdSection) []string {
	accountNumbersSet := make(map[string]bool)
	h.collectAccountNumbersRecursive(sections, accountNumbersSet)

	// Convert set to slice
	var accountNumbers []string
	for accountNumber := range accountNumbersSet {
		accountNumbers = append(accountNumbers, accountNumber)
	}
	return accountNumbers
}

func (h UpdateReportHandler) collectAccountNumbersRecursive(sections []UpdateReportCmdSection, accountNumbersSet map[string]bool) {
	for i := range sections {
		// Collect from items
		for j := range sections[i].Items {
			item := &sections[i].Items[j]
			for k := range item.Formulas {
				accountNumbersSet[item.Formulas[k].AccountNumber] = true
			}
		}
		// Recursively collect from nested sections
		if len(sections[i].Sections) > 0 {
			h.collectAccountNumbersRecursive(sections[i].Sections, accountNumbersSet)
		}
	}
}

// updateFormulaAccountIds recursively updates formula account IDs in all sections
func (h UpdateReportHandler) updateFormulaAccountIds(sections []UpdateReportCmdSection, accountIds map[string]uuid.UUID) error {
	for i := range sections {
		// Update items in this section
		for j := range sections[i].Items {
			item := &sections[i].Items[j]
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
		// Recursively update nested sections
		if len(sections[i].Sections) > 0 {
			if err := h.updateFormulaAccountIds(sections[i].Sections, accountIds); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h UpdateReportHandler) cmdToParams(cmd UpdateReportCmd) report.UpdateReportParams {
	sections := h.convertSections(cmd.Sections)

	return report.UpdateReportParams{
		Title:       cmd.Title,
		AmountTypes: cmd.AmountTypes,
		Sections:    sections,
	}
}

// convertSections recursively converts command sections to domain params
func (h UpdateReportHandler) convertSections(sectionsCmd []UpdateReportCmdSection) []report.UpdateReportParamsSection {
	var sections []report.UpdateReportParamsSection
	for _, sectionCmd := range sectionsCmd {
		var items []report.UpdateReportParamsItem
		for _, itemCmd := range sectionCmd.Items {
			var formulas []report.UpdateReportParamsFormula
			for _, formulaCmd := range itemCmd.Formulas {
				formulas = append(formulas, report.UpdateReportParamsFormula{
					FormulaId: formulaCmd.FormulaId,
					SumFactor: formulaCmd.SumFactor,
					AccountId: formulaCmd.AccountId,
					Rule:      formulaCmd.Rule,
				})
			}

			items = append(items, report.UpdateReportParamsItem{
				ItemId:           itemCmd.ItemId,
				Text:             itemCmd.Text,
				Level:            itemCmd.Level,
				SumFactor:        itemCmd.SumFactor,
				DisplaySumFactor: itemCmd.DisplaySumFactor,
				DataSource:       itemCmd.DataSource,
				Formulas:         formulas,
				IsBreakdownItem:  itemCmd.IsBreakdownItem,
				IsAbleToAddChild: itemCmd.IsAbleToAddChild,
			})
		}

		// Recursively convert nested sections
		var nestedSections []report.UpdateReportParamsSection
		if len(sectionCmd.Sections) > 0 {
			nestedSections = h.convertSections(sectionCmd.Sections)
		}

		sections = append(sections, report.UpdateReportParamsSection{
			SectionId: sectionCmd.SectionId,
			Title:     sectionCmd.Title,
			Items:     items,
			Sections:  nestedSections,
		})
	}

	return sections
}
