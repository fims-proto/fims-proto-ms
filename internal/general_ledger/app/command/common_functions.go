package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"
)

// prepareLineItems prepares line item domain objects and performs necessary checks
func prepareLineItems(
	ctx context.Context,
	readModel query.GeneralLedgerReadModel,
	sobId uuid.UUID,
	commands []LineItemCmd,
) ([]voucher.LineItem, error) {
	// validate account numbers
	var accountNumbers []string
	for _, item := range commands {
		accountNumbers = append(accountNumbers, item.AccountNumber)
	}

	accounts, err := readModel.AccountsByNumbers(ctx, sobId, accountNumbers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to validate account numbers")
	}

	accountIds := utils.SliceToMap(
		accounts,
		func(a query.Account) string { return a.AccountNumber },
		func(a query.Account) uuid.UUID { return a.Id },
	)

	for _, number := range accountNumbers {
		if _, ok := accountIds[number]; !ok {
			return nil, commonErrors.NewSlugError("account-invalidAccountNumber", number)
		}
	}

	// prepare line items
	var lineItems []voucher.LineItem
	for _, item := range commands {
		itemId := item.Id
		if itemId == uuid.Nil {
			itemId = uuid.New()
		}
		lineItem, err := voucher.NewLineItem(
			itemId,
			accountIds[item.AccountNumber],
			item.Text,
			item.Debit,
			item.Credit,
		)
		if err != nil {
			return nil, err
		}
		lineItems = append(lineItems, *lineItem)
	}

	return lineItems, nil
}

// readOrCreatePeriodForVoucher tries to get period id by given transaction time of a voucher, and will also check if the period is closed.
// if no period exists for given transaction time, it creates one
func readOrCreatePeriodForVoucher(
	ctx context.Context,
	repo domain.Repository,
	readModel query.GeneralLedgerReadModel,
	numberingService service.NumberingService,
	sobId uuid.UUID,
	transactionTime time.Time,
) (uuid.UUID, error) {
	fiscalYear := transactionTime.Year()
	periodNumber := int(transactionTime.Month())

	existedPeriodId, err := createPeriodIfNotExists(ctx, repo, readModel, numberingService, sobId, fiscalYear, periodNumber)
	if err != nil {
		return uuid.Nil, err
	}

	p, _ := readModel.PeriodById(ctx, existedPeriodId)
	if p.IsClosed {
		return uuid.Nil, commonErrors.NewSlugError("voucher-periodClosed")
	}

	return existedPeriodId, nil
}

func createPeriodIfNotExists(
	ctx context.Context,
	repo domain.Repository,
	readModel query.GeneralLedgerReadModel,
	numberingService service.NumberingService,
	sobId uuid.UUID,
	fiscalYear int,
	periodNumber int,
) (uuid.UUID, error) {
	existedPeriod, err := readModel.PeriodByFiscalYearAndNumber(ctx, sobId, fiscalYear, periodNumber)
	if err == nil {
		// period exists
		return existedPeriod.Id, nil
	}

	if err.Error() != "period-notFound" {
		return uuid.Nil, errors.Wrap(err, "failed to read period by transaction time")
	}

	// create period
	newPeriodId := uuid.New()
	createPeriodCmd := CreatePeriodCmd{
		SobId:      sobId,
		PeriodId:   newPeriodId,
		FiscalYear: fiscalYear,
		Number:     periodNumber,
	}
	if err = createPeriod(ctx, createPeriodCmd, repo, numberingService); err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to create period")
	}

	return newPeriodId, nil
}
