package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"
	"github/fims-proto/fims-proto-ms/internal/journal/domain/line_item"
	"time"

	"github/fims-proto/fims-proto-ms/internal/journal/app/service"

	"github/fims-proto/fims-proto-ms/internal/journal/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UpdateJournalEntryCmd struct {
	EntryId         uuid.UUID
	LineItems       []LineItemCmd
	TransactionTime time.Time
	User            uuid.UUID
}

type UpdateJournalEntryHandler struct {
	repo           domain.Repository
	accountService service.AccountService
}

func NewUpdateJournalEntryHandler(repo domain.Repository, accountService service.AccountService) UpdateJournalEntryHandler {
	if repo == nil {
		panic("nil repo")
	}

	if accountService == nil {
		panic("nil account service")
	}

	return UpdateJournalEntryHandler{
		repo:           repo,
		accountService: accountService,
	}
}

func (h UpdateJournalEntryHandler) Handle(ctx context.Context, cmd UpdateJournalEntryCmd) error {
	return h.repo.UpdateJournalEntry(
		ctx,
		cmd.EntryId,
		func(j *journal_entry.JournalEntry) (*journal_entry.JournalEntry, error) {
			if len(cmd.LineItems) > 0 {
				// validate account numbers
				var accountNumbers []string
				for _, item := range cmd.LineItems {
					accountNumbers = append(accountNumbers, item.AccountNumber)
				}

				accountIds, err := h.accountService.ValidateExistenceAndGetId(ctx, j.SobId(), accountNumbers)
				if err != nil {
					return nil, errors.Wrap(err, "unable to validate account numbers")
				}

				// prepare line items
				var lineItems []line_item.LineItem
				for _, item := range cmd.LineItems {
					accountId, ok := accountIds[item.AccountNumber]
					if !ok {
						return nil, errors.Errorf("unable to find account id by number %s", item.AccountNumber)
					}

					itemId := item.ItemId
					if itemId == uuid.Nil {
						itemId = uuid.New()
					}
					lineItem, err := line_item.New(
						itemId,
						accountId,
						item.Text,
						item.Debit,
						item.Credit,
					)
					if err != nil {
						return nil, err
					}
					lineItems = append(lineItems, *lineItem)
				}

				if err := j.UpdateLineItems(lineItems, cmd.User); err != nil {
					return nil, err
				}
			}

			if !cmd.TransactionTime.IsZero() {
				period, err := h.accountService.ReadPeriodByTime(ctx, j.SobId(), cmd.TransactionTime)
				if err != nil {
					return nil, errors.Wrap(err, "failed to read period by transaction time")
				}

				if period.IsClosed {
					return nil, errors.New("period is closed")
				}

				if err := j.UpdateTransactionTime(cmd.TransactionTime, period.PeriodId, cmd.User); err != nil {
					return nil, err
				}
			}

			return j, nil
		},
	)
}
