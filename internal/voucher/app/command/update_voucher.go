package command

import (
	"context"
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

func (h UpdateVoucherHandler) Handle(ctx context.Context, cmd UpdateVoucherCmd) error {
	var accNumbers []string
	var lineItems []domain.LineItem
	for _, item := range cmd.LineItems {
		lineItem, err := domain.NewLineItem(
			item.Summary,
			item.AccountNumber,
			item.Debit,
			item.Credit,
		)
		if err != nil {
			return err
		}
		lineItems = append(lineItems, *lineItem)
		accNumbers = append(accNumbers, item.AccountNumber)
	}

	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			if err := h.accountService.ValidateExistence(ctx, accNumbers); err != nil {
				return nil, errors.Wrap(err, "unable to validate account numbers")
			}
			if err := v.Update(lineItems); err != nil {
				return nil, err
			}
			return v, nil
		},
	)
}
