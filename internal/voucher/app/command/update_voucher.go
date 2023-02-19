package command

import (
	"context"
	"time"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain/line_item"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/voucher/app/service"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UpdateVoucherCmd struct {
	VoucherId       uuid.UUID
	HeaderText      string
	LineItems       []LineItemCmd
	TransactionTime time.Time
	Updater         uuid.UUID
}

type UpdateVoucherHandler struct {
	repo           domain.Repository
	accountService service.AccountService
}

func NewUpdateVoucherHandler(repo domain.Repository, accountService service.AccountService) UpdateVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}

	if accountService == nil {
		panic("nil account service")
	}

	return UpdateVoucherHandler{
		repo:           repo,
		accountService: accountService,
	}
}

func (h UpdateVoucherHandler) Handle(ctx context.Context, cmd UpdateVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherId,
		func(j *voucher.Voucher) (*voucher.Voucher, error) {
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

					itemId := item.Id
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

				if err := j.UpdateLineItems(lineItems, cmd.Updater); err != nil {
					return nil, err
				}
			}

			if !cmd.TransactionTime.IsZero() {
				period, err := h.accountService.ReadOrCreatePeriodByTime(ctx, j.SobId(), cmd.TransactionTime)
				if err != nil {
					return nil, errors.Wrap(err, "failed to read period by transaction time")
				}

				if period.IsClosed {
					return nil, commonErrors.NewSlugError("voucher-periodClosed")
				}

				if err := j.UpdateTransactionTime(cmd.TransactionTime, period.Id, cmd.Updater); err != nil {
					return nil, err
				}
			}

			if cmd.HeaderText != "" {
				if err := j.UpdateHeaderText(cmd.HeaderText, cmd.Updater); err != nil {
					return nil, err
				}
			}

			return j, nil
		},
	)
}
