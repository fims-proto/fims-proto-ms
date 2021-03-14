package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RecordVoucherCmd struct {
	Number             string
	AttachmentQuantity uint
	LineItems          []LineItemCmd
	Creator            string
}

type RecordVoucherHandler struct {
	repo           domain.Repository
	accountService AccountService
}

func NewRecordVoucherHandler(repo domain.Repository, accountService AccountService) RecordVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	if accountService == nil {
		panic("nil account service")
	}
	return RecordVoucherHandler{
		repo:           repo,
		accountService: accountService,
	}
}

func (h RecordVoucherHandler) Handle(ctx context.Context, cmd RecordVoucherCmd) (uuid.UUID, error) {
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
			return uuid.Nil, err
		}
		lineItems = append(lineItems, *lineItem)
		accNumbers = append(accNumbers, item.AccountNumber)
	}

	newVoucher, err := domain.NewVoucher(
		uuid.New(),
		cmd.Number,
		time.Now(),
		cmd.AttachmentQuantity,
		lineItems,
		cmd.Creator,
	)
	if err != nil {
		return uuid.Nil, err
	}

	if err = h.accountService.ValidateExistence(ctx, accNumbers); err != nil {
		return uuid.Nil, errors.Wrap(err, "unable to validate account numbers")
	}

	return h.repo.AddVoucher(ctx, newVoucher)
}
