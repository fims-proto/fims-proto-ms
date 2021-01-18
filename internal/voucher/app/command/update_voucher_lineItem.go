package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
)

// UpdateVoucherLineItem ...
type UpdateVoucherLineItemCmd struct{
	VoucherUUID string
	ItemIndex int
	NewItem LineItemCmd
}

type UpdateVoucherLineItemHandler struct {
	repo voucher.Repository
}

func NewUpdateVoucherLineItemHandler(repo  voucher.Repository) UpdateVoucherLineItemHandler {
	if repo == nil{
		panic("nil repo")
	}
	return UpdateVoucherLineItemHandler{repo:repo}

}

func (h UpdateVoucherLineItemHandler) Handle(ctx context.Context, cmd *UpdateVoucherLineItemCmd) error {
	err := h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *voucher.Voucher) (*voucher.Voucher, error){
			if (len(v.LineItems()) <= cmd.ItemIndex){
				panic("LineItem index out of range")
			}
			item, err:= lineitem.NewLineItem(
				cmd.NewItem.Summary,
				cmd.NewItem.AccountNumber,
				cmd.NewItem.Debit,
				cmd.NewItem.Credit,
			)
			if(err==nil){
				// succeeded in NewLineItem
				err:=v.UpdateLineItem(cmd.ItemIndex, item)
				return v, err
			}
			return nil,err
		},
	)
	return err
}