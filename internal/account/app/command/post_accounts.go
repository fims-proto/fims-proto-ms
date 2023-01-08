package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/account/domain/ledger"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
)

type PostAccountsCmd struct {
	PeriodId uuid.UUID
	Records  []PostAccountsRecordCmd
}

type PostAccountsRecordCmd struct {
	AccountId uuid.UUID
	Debit     decimal.Decimal
	Credit    decimal.Decimal
}

type PostAccountsHandler struct {
	repo      domain.Repository
	readModel query.AccountReadModel
}

func NewPostAccountsHandler(
	repo domain.Repository,
	readModel query.AccountReadModel,
) PostAccountsHandler {
	if repo == nil {
		panic("nil account repo")
	}

	if readModel == nil {
		panic("nil read model")
	}

	return PostAccountsHandler{
		repo:      repo,
		readModel: readModel,
	}
}

func (h PostAccountsHandler) Handle(ctx context.Context, cmd PostAccountsCmd) error {
	// all involved accounts
	// prepare all ids and map
	var accountIds []uuid.UUID
	accountsMap := make(map[uuid.UUID]PostAccountsRecordCmd)
	for _, record := range cmd.Records {
		accountIds = append(accountIds, record.AccountId)
		accountsMap[record.AccountId] = record

		//  read all superior accounts
		superiorAccountConfigs, err := h.readModel.SuperiorAccounts(ctx, record.AccountId)
		if err != nil {
			return errors.Wrap(err, "failed to post accounts, cannot read superior accounts")
		}

		for _, superiorAccount := range superiorAccountConfigs {
			superiorRecord := PostAccountsRecordCmd{
				AccountId: superiorAccount.Id,
				Debit:     record.Debit,
				Credit:    record.Credit,
			}
			accountIds = append(accountIds, superiorRecord.AccountId)
			accountsMap[superiorRecord.AccountId] = superiorRecord
		}
	}

	return h.repo.UpdateLedgersByPeriodAndAccountIds(
		ctx,
		cmd.PeriodId,
		accountIds,
		func(accounts []*ledger.Ledger) ([]*ledger.Ledger, error) {
			for _, domainAccount := range accounts {
				record, ok := accountsMap[domainAccount.AccountId()]
				if !ok {
					return nil, errors.Errorf("should not happen, failed to find account %s in accountsMap", domainAccount.AccountId())
				}

				domainAccount.UpdateBalance(record.Debit, record.Credit)
			}

			return accounts, nil
		},
	)
}
