package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github.com/google/uuid"
)

// prepareLineItems prepares line item domain objects and performs necessary checks
func prepareLineItems(
	ctx context.Context,
	repo domain.Repository,
	sobId uuid.UUID,
	commands []LineItemCmd,
) ([]*voucher.LineItem, error) {
	var accountNumbers []string
	var auxiliaryPair []auxiliary_account.AuxiliaryPair
	for _, item := range commands {
		accountNumbers = append(accountNumbers, item.AccountNumber)
		for _, pair := range item.AuxiliaryAccounts {
			auxiliaryPair = append(auxiliaryPair, auxiliary_account.AuxiliaryPair{
				CategoryKey: pair.CategoryKey,
				AccountKey:  pair.AccountKey,
			})
		}
	}

	// validate account numbers
	accounts, err := repo.ReadAccountsByNumbers(ctx, sobId, accountNumbers)
	if err != nil {
		return nil, fmt.Errorf("failed to read accounts: %w", err)
	}

	accountsMap := utils.SliceToMap(
		accounts,
		func(a *account.Account) string { return a.AccountNumber() },
		func(a *account.Account) *account.Account { return a },
	)

	// validate auxiliary account keys
	auxiliaryAccounts, err := repo.ReadAuxiliaryAccountsByPairs(ctx, sobId, auxiliaryPair)
	if err != nil {
		return nil, fmt.Errorf("failed to read auxiliary accounts: %w", err)
	}

	auxiliaryAccountsMap := utils.SliceToMap(
		auxiliaryAccounts,
		func(a *auxiliary_account.AuxiliaryAccount) string { return a.Category().Key() + a.Key() },
		func(a *auxiliary_account.AuxiliaryAccount) *auxiliary_account.AuxiliaryAccount { return a },
	)

	for _, key := range auxiliaryPair {
		if _, ok := auxiliaryAccountsMap[key.CategoryKey+key.AccountKey]; !ok {
			return nil, commonErrors.ErrInvalidAuxiliaryAccountKey(key.CategoryKey, key.AccountKey)
		}
	}

	// prepare line items
	var lineItems []*voucher.LineItem
	for _, item := range commands {
		itemId := item.Id
		if itemId == uuid.Nil {
			itemId = uuid.New()
		}
		a := accountsMap[item.AccountNumber]
		var auxiliaryAccountsForItem []*auxiliary_account.AuxiliaryAccount
		for _, key := range item.AuxiliaryAccounts {
			auxiliaryAccountsForItem = append(auxiliaryAccountsForItem, auxiliaryAccountsMap[key.CategoryKey+key.AccountKey])
		}
		lineItem, err := voucher.NewLineItem(
			itemId,
			a,
			auxiliaryAccountsForItem,
			item.Text,
			item.Amount,
		)
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, lineItem)
	}

	return lineItems, nil
}

// readPeriodIdAndCheck tries to get period id by given transaction date of a voucher, and will also check if the period is closed.
// if no period exists for given transaction date, it creates one
func readPeriodIdAndCheck(
	ctx context.Context,
	repo domain.Repository,
	numberingService service.NumberingService,
	sobId uuid.UUID,
	transactionDate transaction_date.TransactionDate,
) (*period.Period, error) {
	fiscalYear := transactionDate.Year
	periodNumber := transactionDate.Month

	p, err := createPeriodIfNotExists(ctx, createPeriodCmd{
		SobId:      sobId,
		PeriodId:   uuid.Nil,
		FiscalYear: fiscalYear,
		Number:     periodNumber,
	}, repo, numberingService)
	if err != nil {
		return nil, err
	}

	if p.IsClosed() {
		return nil, commonErrors.ErrPeriodClosed()
	}

	return p, nil
}
