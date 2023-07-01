package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"
)

// prepareLineItems prepares line item domain objects and performs necessary checks
func prepareLineItems(
	ctx context.Context,
	repo domain.Repository,
	sobId uuid.UUID,
	commands []LineItemCmd,
) ([]*voucher.LineItem, error) {
	// validate account numbers
	var accountNumbers []string
	for _, item := range commands {
		accountNumbers = append(accountNumbers, item.AccountNumber)
	}

	accounts, err := repo.ReadAccountsByNumbers(ctx, sobId, accountNumbers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to validate account numbers")
	}

	accountsMap := utils.SliceToMap(
		accounts,
		func(a *account.Account) string { return a.AccountNumber() },
		func(a *account.Account) *account.Account { return a },
	)

	for _, number := range accountNumbers {
		if _, ok := accountsMap[number]; !ok {
			return nil, commonErrors.NewSlugError("account-invalidAccountNumber", number)
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
		lineItem, err := voucher.NewLineItem(
			itemId,
			a.Id(),
			a,
			nil,
			item.Text,
			item.Debit,
			item.Credit,
		)
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, lineItem)
	}

	return lineItems, nil
}

// readPeriodIdAndCheck tries to get period id by given transaction time of a voucher, and will also check if the period is closed.
// if no period exists for given transaction time, it creates one
func readPeriodIdAndCheck(ctx context.Context, repo domain.Repository, numberingService service.NumberingService, sobId uuid.UUID, transactionTime time.Time) (uuid.UUID, error) {
	fiscalYear := transactionTime.Year()
	periodNumber := int(transactionTime.Month())

	p, err := createPeriodIfNotExists(ctx, createPeriodCmd{
		SobId:      sobId,
		PeriodId:   uuid.Nil,
		FiscalYear: fiscalYear,
		Number:     periodNumber,
	}, repo, numberingService)
	if err != nil {
		return uuid.Nil, err
	}

	if p.IsClosed() {
		return uuid.Nil, commonErrors.NewSlugError("voucher-periodClosed")
	}

	return p.Id(), nil
}
