package account

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/account/app/command"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"

	accountPort "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	accountInterface accountPort.AccountInterface
}

func NewIntraProcessAdapter(accountInterface accountPort.AccountInterface) IntraProcessAdapter {
	return IntraProcessAdapter{accountInterface: accountInterface}
}

func (i IntraProcessAdapter) ValidateExistenceAndGetId(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]uuid.UUID, error) {
	accountConfigurations, err := i.accountInterface.ReadAccountConfigurationsByNumbers(ctx, sobId, accountNumbers)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate existence")
	}
	accountConfigsMap := make(map[string]uuid.UUID)
	for _, config := range accountConfigurations {
		accountConfigsMap[config.AccountNumber] = config.Id
	}
	return accountConfigsMap, nil
}

func (i IntraProcessAdapter) ReadAccountConfigurationsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]query.Account, error) {
	accountConfigurations, err := i.accountInterface.ReadAccountConfigurationsByIds(ctx, accountIds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read account configuration")
	}
	accountConfigsMap := make(map[uuid.UUID]query.Account)
	for _, config := range accountConfigurations {
		accountConfigsMap[config.Id] = config
	}
	return accountConfigsMap, nil
}

func (i IntraProcessAdapter) ReadPeriodByTime(ctx context.Context, sobId uuid.UUID, transactionTime time.Time) (query.Period, error) {
	return i.accountInterface.ReadPeriodByTime(ctx, sobId, transactionTime)
}

func (i IntraProcessAdapter) ReadPeriodsByIds(ctx context.Context, periodIds []uuid.UUID) (map[uuid.UUID]query.Period, error) {
	periods, err := i.accountInterface.ReadPeriodsByIds(ctx, periodIds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read periods")
	}
	periodsMap := make(map[uuid.UUID]query.Period)
	for _, period := range periods {
		periodsMap[period.PeriodId] = period
	}
	return periodsMap, nil
}

func (i IntraProcessAdapter) PostJournalEntry(ctx context.Context, entry journal_entry.JournalEntry) error {
	var records []command.PostAccountsRecordCmd
	for _, item := range entry.LineItems() {
		records = append(records, command.PostAccountsRecordCmd{
			AccountId: item.AccountId(),
			Debit:     item.Debit(),
			Credit:    item.Credit(),
		})
	}

	return i.accountInterface.PostAccounts(ctx, command.PostAccountsCmd{
		PeriodId: entry.PeriodId(),
		Records:  records,
	})
}
