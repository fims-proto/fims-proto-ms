package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/user"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RecordVoucherCmd struct {
	Sob                string
	VoucherType        string
	AttachmentQuantity uint
	LineItems          []LineItemCmd
	Creator            string
}

type RecordVoucherHandler struct {
	repo           domain.Repository
	accountService AccountService
	counterService CounterService
}

func NewRecordVoucherHandler(repo domain.Repository, accountService AccountService, counterService CounterService) RecordVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	if accountService == nil {
		panic("nil account service")
	}
	if counterService == nil {
		panic("nil counter service")
	}
	return RecordVoucherHandler{
		repo:           repo,
		accountService: accountService,
		counterService: counterService,
	}
}

func (h RecordVoucherHandler) Handle(ctx context.Context, cmd RecordVoucherCmd) (uuid.UUID, error) {
	if err := user.VerifyAuth("current-user", cmd.Sob, "voucher", "create"); err != nil {
		return uuid.Nil, errors.Wrap(err, "failed verifing permission")
	}

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

	voucherType, err := domain.NewVoucherTypeFromString(cmd.VoucherType)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "unable to use voucher type")
	}

	identifier, err := h.counterService.GetNextIdentifier(ctx, cmd.Sob, voucherType.String())
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "unable to generate next number")
	}

	newVoucher, err := domain.NewVoucher(
		cmd.Sob,
		uuid.New(),
		voucherType,
		identifier,
		time.Now(),
		cmd.AttachmentQuantity,
		lineItems,
		cmd.Creator,
	)
	if err != nil {
		return uuid.Nil, err
	}

	if err = h.accountService.ValidateExistence(ctx, cmd.Sob, accNumbers); err != nil {
		return uuid.Nil, errors.Wrap(err, "unable to validate account numbers")
	}

	return h.repo.AddVoucher(ctx, newVoucher)
}
