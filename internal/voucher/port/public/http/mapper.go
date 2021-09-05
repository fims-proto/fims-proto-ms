package http

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
)

func (r LineItemRequest) mapToCommand() command.LineItemCmd {
	return command.LineItemCmd{
		Summary:       r.Summary,
		AccountNumber: r.AccountNumber,
		Debit:         r.Debit,
		Credit:        r.Credit,
	}
}

func (r RecordVoucherRequest) mapToCommand() command.RecordVoucherCmd {
	itemCmd := []command.LineItemCmd{}
	for _, item := range r.LineItems {
		itemCmd = append(itemCmd, item.mapToCommand())
	}
	return command.RecordVoucherCmd{
		VoucherType:        r.VoucherType,
		AttachmentQuantity: uint(r.AttachmentQuantity),
		LineItems:          itemCmd,
		Creator:            r.Creator,
	}
}

func mapFromLineItemQuery(q query.LineItem) LineItemResponse {
	return LineItemResponse{
		Id:            q.Id.String(),
		Summary:       q.Summary,
		AccountNumber: q.AccountNumber,
		Debit:         q.Debit,
		Credit:        q.Credit,
	}
}

func mapFromVoucherQuery(q query.Voucher) VoucherResponse {
	itemRes := []LineItemResponse{}
	for _, item := range q.LineItems {
		itemRes = append(itemRes, mapFromLineItemQuery(item))
	}
	return VoucherResponse{
		Sob:                q.Sob,
		Id:                 q.Id.String(),
		Type:               q.VoucherType,
		Number:             string(q.Number),
		CreatedAt:          q.CreatedAt,
		AttachmentQuantity: int(q.AttachmentQuantity),
		LineItems:          itemRes,
		Debit:              q.Debit,
		Credit:             q.Credit,
		Creator:            q.Creator,
		Reviewer:           q.Reviewer,
		Auditor:            q.Auditor,
		IsReviewed:         q.IsReviewed,
		IsAudited:          q.IsAudited,
		IsPosted:           q.IsPosted,
	}
}
