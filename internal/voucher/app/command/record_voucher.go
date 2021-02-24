package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RecordVoucherCmd struct {
	Number             string
	CreatedAt          time.Time
	AttachmentQuantity uint
	LineItems          []LineItemCmd
	CreatorUUID        string
}

type RecordVoucherHandler struct {
	repo       voucher.Repository
	accService AccountService
}

func NewRecordVoucherHandler(repo voucher.Repository, accService AccountService) RecordVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	if accService == nil {
		panic("nil account service")
	}
	return RecordVoucherHandler{
		repo:       repo,
		accService: accService,
	}
}

func (h RecordVoucherHandler) Handle(ctx context.Context, cmd RecordVoucherCmd) (uuid.UUID, error) {
	// object conversion, outside in: LineItemCmd -> domain/LineItem
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
			return uuid.Nil, err
		}
		lineItems = append(lineItems, *lineItem)
		accNumbers = append(accNumbers, item.AccountNumber)
	}

	// object conversion, outside in: VoucherCmd -> domain/Voucher
	newVoucher, err := voucher.NewVoucher(
		uuid.New(),
		cmd.Number,
		cmd.CreatedAt,
		cmd.AttachmentQuantity,
		lineItems,
		cmd.CreatorUUID,
	)
	if err != nil {
		return uuid.Nil, err
	}

	if err = h.accService.ValidateExistence(ctx, accNumbers); err != nil {
		return uuid.Nil, errors.Wrap(err, "unable to validate account numbers")
	}

	return h.repo.AddVoucher(ctx, newVoucher)
}
