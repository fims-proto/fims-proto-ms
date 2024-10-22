package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"
)

type InitializeCmd struct {
	SobId uuid.UUID
}

type InitializeHandler struct {
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService

	accounts map[string]uuid.UUID // key: accountNumber, value: accountId
}

func NewInitializeHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) InitializeHandler {
	if repo == nil {
		panic("nil repo")
	}

	if generalLedgerService == nil {
		panic("nil general ledger service")
	}

	return InitializeHandler{
		repo:                 repo,
		generalLedgerService: generalLedgerService,
		accounts:             make(map[string]uuid.UUID),
	}
}

func (h *InitializeHandler) Handle(ctx context.Context, cmd InitializeCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.handle(txCtx, cmd)
	})
}

func (h *InitializeHandler) handle(ctx context.Context, cmd InitializeCmd) error {
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %w", err)
	}

	var balanceSheetCmd InitializeCmdReport
	var incomeStatementCmd InitializeCmdReport

	balanceSheetFile, err := os.ReadFile(filepath.Join(workDir, "dataload", "xqykjzz", "report_balance_sheet.json"))
	if err != nil {
		return fmt.Errorf("could not read balance sheet json file: %w", err)
	}

	incomeStatementFile, err := os.ReadFile(filepath.Join(workDir, "dataload", "xqykjzz", "report_income_statement.json"))
	if err != nil {
		return fmt.Errorf("could not read income statement json file: %w", err)
	}

	if err = json.Unmarshal(balanceSheetFile, &balanceSheetCmd); err != nil {
		return fmt.Errorf("could not unmarshal balance sheet json: %w", err)
	}
	if err = json.Unmarshal(incomeStatementFile, &incomeStatementCmd); err != nil {
		return fmt.Errorf("could not unmarshal income statement json: %w", err)
	}

	if err = h.prepareAccounts(ctx, cmd.SobId, balanceSheetCmd, incomeStatementCmd); err != nil {
		return fmt.Errorf("could not prepare accounts with account numbers: %w", err)
	}

	balanceSheetReport, err := h.convertReport(cmd.SobId, balanceSheetCmd)
	if err != nil {
		return fmt.Errorf("could not convert balance sheet report: %w", err)
	}
	incomeStatementReport, err := h.convertReport(cmd.SobId, incomeStatementCmd)
	if err != nil {
		return fmt.Errorf("could not convert income statement report: %w", err)
	}

	if err = h.repo.CreateReports(ctx, []*report.Report{balanceSheetReport, incomeStatementReport}); err != nil {
		return fmt.Errorf("could not create reports: %w", err)
	}

	return nil
}

func (h *InitializeHandler) convertReport(sobId uuid.UUID, cmd InitializeCmdReport) (*report.Report, error) {
	var sections []*report.Section
	for _, cmdSection := range cmd.Sections {
		section, err := h.convertSection(cmdSection)
		if err != nil {
			return nil, err
		}
		sections = append(sections, section)
	}

	return report.New(
		uuid.New(),
		sobId,
		uuid.Nil,
		cmd.Title,
		true,
		cmd.Class,
		cmd.AmountTypes,
		sections,
	)
}

func (h *InitializeHandler) convertSection(cmd InitializeCmdSection) (*report.Section, error) {
	var subSections []*report.Section
	for _, cmdSubSection := range cmd.Sections {
		subSection, err := h.convertSection(cmdSubSection)
		if err != nil {
			return nil, err
		}
		subSections = append(subSections, subSection)
	}

	var items []*report.Item
	for _, cmdItem := range cmd.Items {
		item, err := h.convertItem(cmdItem)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return report.NewSection(
		uuid.New(),
		cmd.Title,
		nil,
		subSections,
		items,
	)
}

func (h *InitializeHandler) convertItem(cmd InitializeCmdItem) (*report.Item, error) {
	var formulas []*report.Formula
	for _, cmdFormula := range cmd.Formulas {
		formula, err := h.convertFormula(cmdFormula)
		if err != nil {
			return nil, err
		}
		formulas = append(formulas, formula)
	}

	return report.NewItem(
		uuid.New(),
		cmd.Text,
		cmd.Level,
		cmd.SumFactor,
		cmd.DisplaySumFactor,
		cmd.DataSource,
		formulas,
		nil,
		cmd.IsBreakdownItem,
		cmd.IsDeletable,
		cmd.IsTextModifiable,
		cmd.IsDraggable,
		cmd.IsAbleToAddChild,
		cmd.IsAbleToAddLeaf,
	)
}

func (h *InitializeHandler) convertFormula(cmd InitializeCmdFormula) (*report.Formula, error) {
	accountId, ok := h.accounts[cmd.AccountNumber]
	if !ok {
		return nil, fmt.Errorf("could not find account number %s", cmd.AccountNumber)
	}

	return report.NewFormula(
		uuid.New(),
		accountId,
		cmd.SumFactor,
		cmd.Rule,
		nil,
	)
}

func (h *InitializeHandler) prepareAccounts(ctx context.Context, sobId uuid.UUID, cmds ...InitializeCmdReport) error {
	var accountNumbers []string

	for _, cmd := range cmds {
		for _, section := range cmd.Sections {
			accountNumbers = append(accountNumbers, collectAccountNumbersFromSection(section)...)
		}
	}

	accountIds, err := h.generalLedgerService.ReadAccountIdsByNumbers(ctx, sobId, accountNumbers)
	if err != nil {
		return err
	}

	h.accounts = accountIds
	return nil
}

func collectAccountNumbersFromSection(section InitializeCmdSection) []string {
	var result []string
	for _, item := range section.Items {
		for _, formula := range item.Formulas {
			result = append(result, formula.AccountNumber)
		}
	}

	for _, subSection := range section.Sections {
		result = append(result, collectAccountNumbersFromSection(subSection)...)
	}

	return result
}
