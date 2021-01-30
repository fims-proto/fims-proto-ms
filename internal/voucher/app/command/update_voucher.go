package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github.com/pkg/errors"
)

type UpdateVoucherCmd struct {
	VoucherUUID string
	LineItems   []LineItemCmd
}

type UpdateVoucherHandler struct {
	repo       voucher.Repository
	accService AccountService
}

func NewUpdateVoucherHandler(repo voucher.Repository, accService AccountService) UpdateVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	if accService == nil {
		panic("nil account service")
	}
	return UpdateVoucherHandler{
		repo:       repo,
		accService: accService,
	}
}

func (h UpdateVoucherHandler) Handle(ctx context.Context, cmd UpdateVoucherCmd) error {
	var accNumbers []string
	var lineItems []lineitem.LineItem
	for _, item := range cmd.LineItems {
		lineItem, err := lineitem.NewLineItem(
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
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			if err := h.accService.ValidateExistence(ctx, accNumbers); err != nil {
				return nil, errors.Wrap(err, "unable to validate account numbers")
			}
			if err := v.Update(lineItems); err != nil {
				return nil, err
			}
			return v, nil
		},
	)
}
