package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UpdateVoucherCmd struct {
	VoucherUUID     uuid.UUID
	LineItems       []LineItemCmd
	TransactionTime time.Time
}

type UpdateVoucherHandler struct {
	repo           domain.Repository
	accountService AccountService
}

func NewUpdateVoucherHandler(repo domain.Repository, accountService AccountService) UpdateVoucherHandler {
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

func (h UpdateVoucherHandler) Handle(ctx context.Context, cmd UpdateVoucherCmd) (err error) {
	log.Info(ctx, "handle updating voucher")
	log.Debug(ctx, "handle updating voucher, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle updating failed")
		}
	}()

	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			if len(cmd.LineItems) > 0 {
				// validate account numbers
				log.Info(ctx, "validating line items")
				var accountNumbers []string
				for _, item := range cmd.LineItems {
					accountNumbers = append(accountNumbers, item.AccountNumber)
				}
				accountIds, err := h.accountService.ValidateExistenceAndGetId(ctx, v.SobId(), accountNumbers)
				if err != nil {
					return nil, errors.Wrap(err, "unable to validate account numbers")
				}

				// prepare line items
				var lineItems []*domain.LineItem
				for _, item := range cmd.LineItems {
					accountId, ok := accountIds[item.AccountNumber]
					if !ok {
						return nil, errors.Wrapf(err, "unable to find account id by number %s", item.AccountNumber)
					}
					itemId := item.Id
					if itemId == uuid.Nil {
						itemId = uuid.New()
					}
					lineItem, err := domain.NewLineItem(
						itemId,
						accountId,
						item.Summary,
						item.Debit,
						item.Credit,
					)
					if err != nil {
						return nil, err
					}
					lineItems = append(lineItems, lineItem)
				}

				log.Info(ctx, "updating voucher line items")
				if err := v.UpdateLineItems(lineItems); err != nil {
					return nil, err
				}
			}

			if !cmd.TransactionTime.IsZero() {
				log.Info(ctx, "updating voucher transaction time")
				if err := v.UpdateTransactionTime(cmd.TransactionTime); err != nil {
					return nil, err
				}
			}

			return v, nil
		},
	)
}
