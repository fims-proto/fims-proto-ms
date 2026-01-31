package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"

	"github.com/google/uuid"
)

type AddItemCmd struct {
	SobId               uuid.UUID
	ReportId            uuid.UUID
	SectionId           uuid.UUID
	InsertAfterSequence int
	Text                string
	Level               int
	SumFactor           int
	DataSource          string
	Formulas            []FormulaCmd
	IsBreakdownItem     bool
	IsAbleToAddChild    bool
}

type AddItemHandler struct {
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService
}

func NewAddItemHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) AddItemHandler {
	if repo == nil {
		panic("nil repo")
	}

	if generalLedgerService == nil {
		panic("nil general ledger service")
	}

	return AddItemHandler{
		repo:                 repo,
		generalLedgerService: generalLedgerService,
	}
}

func (h AddItemHandler) Handle(ctx context.Context, cmd AddItemCmd) (uuid.UUID, error) {
	var itemId uuid.UUID

	err := h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		var txErr error
		itemId, txErr = h.addItem(txCtx, cmd)
		return txErr
	})

	return itemId, err
}

func (h AddItemHandler) addItem(ctx context.Context, cmd AddItemCmd) (uuid.UUID, error) {
	itemId := uuid.New()

	err := h.repo.UpdateReport(ctx, cmd.ReportId, func(r *report.Report) (*report.Report, error) {
		// Prepare account ids for formulas
		var accountNumbers []string
		for _, formulaCmd := range cmd.Formulas {
			accountNumbers = append(accountNumbers, formulaCmd.AccountNumber)
		}
		accountIds, err := h.generalLedgerService.ReadAccountIdsByNumbers(ctx, cmd.SobId, accountNumbers)
		if err != nil {
			return nil, fmt.Errorf("failed to read account ids by numbers: %w", err)
		}

		// Prepare formulas
		var formulas []*report.Formula
		for index, formulaCmd := range cmd.Formulas {
			accountId, ok := accountIds[formulaCmd.AccountNumber]
			if !ok {
				return nil, fmt.Errorf("failed to find account id by account number: %s", formulaCmd.AccountNumber)
			}
			formula, err := report.NewFormula(uuid.New(), index+1, accountId, formulaCmd.SumFactor, formulaCmd.Rule, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create formula: %w", err)
			}
			formulas = append(formulas, formula)
		}

		// Create new item
		// Note: sequence will be auto-assigned by the Section.AddItem method
		newItem, err := report.NewItem(
			itemId,
			cmd.Text,
			cmd.Level,
			1,  // Temporary sequence, will be renumbered by AddItem
			"", // Keep none
			cmd.SumFactor,
			false, // displaySumFactor - default to false
			cmd.DataSource,
			formulas,
			nil,  // amounts - will be calculated during report generation
			true, // isEditable - new items are editable
			cmd.IsBreakdownItem,
			cmd.IsAbleToAddChild,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create item: %w", err)
		}

		// Add item to section at specified position
		if err := r.AddItemToSection(cmd.SectionId, newItem, cmd.InsertAfterSequence); err != nil {
			return nil, err
		}

		return r, nil
	})

	return itemId, err
}
