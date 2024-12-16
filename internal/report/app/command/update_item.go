package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/data_source"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"
)

type UpdateItemCmd struct {
	SobId      uuid.UUID
	Id         uuid.UUID
	Text       string
	SumFactor  *int // zero value of int will make confusion
	DataSource string
	Formulas   []FormulaCmd
}

type FormulaCmd struct {
	SumFactor     int
	AccountNumber string
	Rule          string
}

type UpdateItemHandler struct {
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService
}

func NewUpdateItemHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) UpdateItemHandler {
	if repo == nil {
		panic("nil repo")
	}
	if generalLedgerService == nil {
		panic("nil general ledger service")
	}
	return UpdateItemHandler{
		repo:                 repo,
		generalLedgerService: generalLedgerService,
	}
}

func (h UpdateItemHandler) Handle(ctx context.Context, cmd UpdateItemCmd) error {
	return h.repo.UpdateItem(ctx, cmd.Id, func(i *report.Item) (*report.Item, error) {
		// update text
		if err := i.UpdateText(cmd.Text); err != nil {
			return nil, err
		}

		// update sum factor
		if cmd.SumFactor != nil {
			if err := i.UpdateSumFactor(*cmd.SumFactor); err != nil {
				return nil, err
			}
		}

		// update data source and formulas
		// prepare account ids
		var accountNumbers []string
		for _, formulaCmd := range cmd.Formulas {
			accountNumbers = append(accountNumbers, formulaCmd.AccountNumber)
		}
		accountIds, err := h.generalLedgerService.ReadAccountIdsByNumbers(ctx, cmd.SobId, accountNumbers)
		if err != nil {
			return nil, fmt.Errorf("failed to read account ids by numbers: %w", err)
		}

		// prepare formulas
		var formulas []*report.Formula
		for _, formulaCmd := range cmd.Formulas {
			accountId, ok := accountIds[formulaCmd.AccountNumber]
			if !ok {
				return nil, fmt.Errorf("failed to find account id by account number: %s", formulaCmd.AccountNumber)
			}
			formula, err := report.NewFormula(uuid.New(), accountId, formulaCmd.SumFactor, formulaCmd.Rule, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create formula: %w", err)
			}
			formulas = append(formulas, formula)
		}

		dataSource, err := data_source.FromString(cmd.DataSource)
		if err != nil {
			return nil, err
		}
		if err = i.UpdateDataSource(dataSource, formulas); err != nil {
			return nil, err
		}

		return i, nil
	})
}
