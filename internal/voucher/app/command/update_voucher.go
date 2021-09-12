package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UpdateVoucherCmd struct {
	VoucherUUID uuid.UUID
	LineItems   []LineItemCmd
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

	var accNumbers []string
	var lineItems []*domain.LineItem
	for _, item := range cmd.LineItems {
		lineItem, err := domain.NewLineItem(
			item.Id,
			item.Summary,
			item.AccountNumber,
			item.Debit,
			item.Credit,
		)
		if err != nil {
			return err
		}
		lineItems = append(lineItems, lineItem)
		accNumbers = append(accNumbers, item.AccountNumber)
	}

	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			log.Info(ctx, "validating account number")
			if err := h.accountService.ValidateExistence(ctx, v.Sob(), accNumbers); err != nil {
				return nil, errors.Wrap(err, "unable to validate account numbers")
			}
			log.Info(ctx, "updating voucher")
			if err := v.Update(lineItems); err != nil {
				return nil, err
			}
			return v, nil
		},
	)
}
