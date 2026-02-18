package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InitializeLedgersBalanceCmd struct {
	SobId   uuid.UUID
	Ledgers []InitializeLedgersBalanceItemCmd
}

type InitializeLedgersBalanceItemCmd struct {
	AccountNumber  string
	OpeningBalance decimal.Decimal
	validated      bool // for command validation
}

type InitializeLedgersBalanceHandler struct {
	repo       domain.Repository
	sobService service.SobService
}

func NewInitializeLedgersBalanceHandler(repo domain.Repository, sobService service.SobService) InitializeLedgersBalanceHandler {
	if repo == nil {
		panic("nil repo")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	return InitializeLedgersBalanceHandler{
		repo:       repo,
		sobService: sobService,
	}
}

func (h InitializeLedgersBalanceHandler) Handle(ctx context.Context, cmd InitializeLedgersBalanceCmd) error {
	// check if first period is closed
	firstPeriod, err := h.repo.ReadFirstPeriod(ctx, cmd.SobId)
	if err != nil {
		return fmt.Errorf("failed to read first period: %w", err)
	}
	if firstPeriod.IsClosed() {
		return errors.ErrPeriodClosed()
	}

	// prepare ledgers to be updated
	type ledgerRecord struct {
		accountId      uuid.UUID
		openingBalance decimal.Decimal
	}

	subAccountsWithSuperiors, err := h.repo.ReadAllSubAccountsWithSuperiors(ctx, cmd.SobId)
	if err != nil {
		return fmt.Errorf("failed to read all sub accounts: %w", err)
	}
	cmdMap := utils.SliceToMap(cmd.Ledgers, func(l InitializeLedgersBalanceItemCmd) string {
		return l.AccountNumber
	}, func(l InitializeLedgersBalanceItemCmd) *InitializeLedgersBalanceItemCmd {
		return &l
	})

	var ledgerRecords []ledgerRecord
	for _, subAccount := range subAccountsWithSuperiors {
		itemCmd, ok := cmdMap[subAccount.AccountNumber()]
		if !ok {
			// means command doesn't provide this account
			continue
		}
		itemCmd.validated = true // mark found
		// for leaves accounts
		ledgerRecords = append(ledgerRecords, ledgerRecord{
			accountId:      subAccount.Id(),
			openingBalance: itemCmd.OpeningBalance,
		})
		// for superiors accounts
		superior := subAccount.SuperiorAccount()
		for superior != nil {
			ledgerRecords = append(ledgerRecords, ledgerRecord{
				accountId:      superior.Id(),
				openingBalance: itemCmd.OpeningBalance,
			})
			superior = superior.SuperiorAccount()
		}
	}
	// if command item remains un-validated, means input gives account we don't know
	for _, itemCmd := range cmdMap {
		if !itemCmd.validated {
			return fmt.Errorf("accept only sub-accounts, but got invalid account: %s", itemCmd.AccountNumber)
		}
	}

	// ledger records could contain duplicated superiors, merge them
	ledgersMap := utils.SliceToMapMerge(ledgerRecords, func(l ledgerRecord) uuid.UUID {
		return l.accountId
	}, func(l ledgerRecord) decimal.Decimal {
		return l.openingBalance
	}, func(existing decimal.Decimal, replacement decimal.Decimal) decimal.Decimal {
		return existing.Add(replacement)
	})

	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.updateAndTrialBalance(txCtx, cmd.SobId, firstPeriod.Id(), ledgersMap)
	})
}

func (h InitializeLedgersBalanceHandler) updateAndTrialBalance(ctx context.Context, sobId, firstPeriodId uuid.UUID, ledgersMap map[uuid.UUID]decimal.Decimal) error {
	if err := h.repo.UpdateLedgersByPeriodAndAccountIds(
		ctx,
		firstPeriodId,
		utils.MapToKeySlice(ledgersMap),
		func(ledgers []*ledger.Ledger) ([]*ledger.Ledger, error) {
			for _, l := range ledgers {
				openingBalance, ok := ledgersMap[l.AccountId()]
				if !ok {
					return nil, fmt.Errorf("should not happen, failed to find account %s in ledgersMap", l.AccountId())
				}

				l.UpdateOpeningBalance(openingBalance)
			}

			return ledgers, nil
		},
	); err != nil {
		return fmt.Errorf("failed to update ledgers: %w", err)
	}

	if err := trialBalance(ctx, h.repo, sobId, firstPeriodId); err != nil {
		return fmt.Errorf("not balance: %w", err)
	}

	return nil
}
