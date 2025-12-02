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
	for i, cmdSection := range cmd.Sections {
		section, err := h.convertSection(cmdSection, i+1)
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

func (h *InitializeHandler) convertSection(cmd InitializeCmdSection, sequence int) (*report.Section, error) {
	var subSections []*report.Section
	for i, cmdSubSection := range cmd.Sections {
		subSection, err := h.convertSection(cmdSubSection, i+1)
		if err != nil {
			return nil, err
		}
		subSections = append(subSections, subSection)
	}

	var items []*report.Item
	for i, cmdItem := range cmd.Items {
		item, err := h.convertItem(cmdItem, i+1)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return report.NewSection(
		uuid.New(),
		cmd.Title,
		sequence,
		nil,
		subSections,
		items,
	)
}

func (h *InitializeHandler) convertItem(cmd InitializeCmdItem, sequence int) (*report.Item, error) {
	var formulas []*report.Formula
	for i, cmdFormula := range cmd.Formulas {
		formula, err := h.convertFormula(cmdFormula, i+1)
		if err != nil {
			return nil, err
		}
		formulas = append(formulas, formula)
	}

	return report.NewItem(
		uuid.New(),
		cmd.Text,
		cmd.Level,
		sequence,
		cmd.SumFactor,
		cmd.DisplaySumFactor,
		cmd.DataSource,
		formulas,
		nil,
		cmd.IsEditable,
		cmd.IsBreakdownItem,
		cmd.IsAbleToAddChild,
		cmd.IsAbleToAddLeaf,
	)
}

func (h *InitializeHandler) convertFormula(cmd InitializeCmdFormula, sequence int) (*report.Formula, error) {
	accountId, ok := h.accounts[cmd.AccountNumber]
	if !ok {
		return nil, fmt.Errorf("could not find account number %s", cmd.AccountNumber)
	}

	return report.NewFormula(
		uuid.New(),
		sequence,
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
