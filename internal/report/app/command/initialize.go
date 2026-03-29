package command

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	appService "github/fims-proto/fims-proto-ms/internal/report/app/service"
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
	sobService           appService.SobService

	accounts    map[string]uuid.UUID // key: rawAccountNumber, value: accountId
	codeLengths []int
}

func NewInitializeHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService, sobService appService.SobService) InitializeHandler {
	if repo == nil {
		panic("nil repo")
	}

	if generalLedgerService == nil {
		panic("nil general ledger service")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	return InitializeHandler{
		repo:                 repo,
		generalLedgerService: generalLedgerService,
		sobService:           sobService,
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

	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.repo.CreateReports(txCtx, []*report.Report{balanceSheetReport, incomeStatementReport})
	})
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
		cmd.SectionType,
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

	return report.NewItem(uuid.New(), cmd.Text, cmd.Level, sequence, cmd.ItemType, cmd.SumFactor, cmd.DisplaySumFactor, cmd.DataSource, formulas, nil, cmd.IsEditable, cmd.IsBreakdownItem, cmd.IsAbleToAddChild)
}

func (h *InitializeHandler) convertFormula(cmd InitializeCmdFormula, sequence int) (*report.Formula, error) {
	// Convert human-readable account number to raw format for lookup
	rawAccountNumber, err := account.RawFromReadable(cmd.AccountNumber, h.codeLengths)
	if err != nil {
		return nil, fmt.Errorf("could not convert account number %s: %w", cmd.AccountNumber, err)
	}

	accountId, ok := h.accounts[rawAccountNumber]
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
	// Fetch SoB to get codeLengths for conversion
	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return fmt.Errorf("could not read sob: %w", err)
	}
	h.codeLengths = sob.AccountsCodeLength

	// Collect all human-readable account numbers and convert to raw format
	var rawAccountNumbers []string
	humanReadableToRaw := make(map[string]string) // human-readable -> raw

	for _, cmd := range cmds {
		for _, section := range cmd.Sections {
			for _, humanNum := range collectAccountNumbersFromSection(section) {
				// Skip if already converted
				if _, exists := humanReadableToRaw[humanNum]; exists {
					continue
				}
				// Convert to raw format
				rawNum, err := account.RawFromReadable(humanNum, h.codeLengths)
				if err != nil {
					return fmt.Errorf("could not convert account number %s: %w", humanNum, err)
				}
				humanReadableToRaw[humanNum] = rawNum
				rawAccountNumbers = append(rawAccountNumbers, rawNum)
			}
		}
	}

	// Query by raw account numbers
	accountIds, err := h.generalLedgerService.ReadAccountIdsByRawNumbers(ctx, sobId, rawAccountNumbers)
	if err != nil {
		return err
	}

	// Store map keyed by raw account numbers
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
