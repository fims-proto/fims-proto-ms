package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"time"
)

type RecordVoucherCmd struct {
	UUID string
	// TODO 字号
	Number             uint
	CreatedAt          time.Time
	AttachmentQuantity uint
	LineItems          []LineItemCmd
	Debit              string
	Credit             string
	CreatorUUID        string
}

// TODO discuss should LineItemCmd or in a types.go file
// type LineItemCmd struct {
// 	Summary       string
// 	AccountNumber string
// 	Debit         string
// 	Credit        string
// }

type RecordVoucherHandler struct {
	repo voucher.ActionModel
}

func NewRecordVoucherHandler(repo voucher.ActionModel) RecordVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}
	return RecordVoucherHandler{repo: repo}
}

func (h RecordVoucherHandler) Handle(ctx context.Context, cmd RecordVoucherCmd) error {
	// object conversion, outside in: LineItemCmd -> domain/LineItem
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
	}

	// object conversion, outside in: VoucherCmd -> domain/Voucher
	newVoucher, err := voucher.NewVoucher(
		cmd.UUID,
		cmd.Number,
		cmd.CreatedAt,
		cmd.AttachmentQuantity,
		lineItems,
		cmd.CreatorUUID,
	)
	if err != nil {
		return err
	}

	return h.repo.AddVoucher(ctx, newVoucher)
}