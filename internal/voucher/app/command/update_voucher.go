package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain/lineitem"
)

type UpdateVoucherCmd struct {
	UUID string
	LineItems []LineItemCmd
}

type UpdateVoucherHandler struct {
	repo voucher.Repository
}

func NewUpdateVoucherHandler(repo  voucher.Repository) UpdateVoucherHandler {
	if repo == nil{
		panic("nil repo")
	}
	return UpdateVoucherHandler{repo:repo}

}

func (h UpdateVoucherHandler) Handle(ctx context.Context, cmd *UpdateVoucherCmd) error {
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

	err := h.repo.UpdateVoucher(
		ctx,
		cmd.UUID,
		func( v *voucher.Voucher) (*voucher.Voucher,error){
			err:= v.Update(lineItems)
			if err != nil{
				return nil, err 
			}
			return v, nil
		},		
	)
	return err
}