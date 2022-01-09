package http

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"

	"github.com/google/uuid"
)

func (r LineItemRequest) mapToCommand() command.LineItemCmd {
	id, _ := uuid.Parse(r.Id)
	return command.LineItemCmd{
		Id:            id,
		Summary:       r.Summary,
		AccountNumber: r.AccountNumber,
		Debit:         r.Debit,
		Credit:        r.Credit,
	}
}

func (r CreateVoucherRequest) mapToCommand() command.CreateVoucherCmd {
	var itemCmd []command.LineItemCmd
	for _, item := range r.LineItems {
		itemCmd = append(itemCmd, item.mapToCommand())
	}
	return command.CreateVoucherCmd{
		VoucherType:        r.VoucherType,
		AttachmentQuantity: uint(r.AttachmentQuantity),
		LineItems:          itemCmd,
		Creator:            r.Creator,
		TransactionTime:    r.TransactionTime,
	}
}

func mapFromLineItemQuery(q query.LineItem) LineItemResponse {
	return LineItemResponse{
		Id:        q.Id.String(),
		AccountId: q.AccountId.String(),
		Summary:   q.Summary,
		Debit:     q.Debit,
		Credit:    q.Credit,
	}
}

func mapFromVoucherQuery(q query.Voucher) VoucherResponse {
	var itemRes []LineItemResponse
	for _, item := range q.LineItems {
		itemRes = append(itemRes, mapFromLineItemQuery(item))
	}
	return VoucherResponse{
		SobId:              q.SobId.String(),
		Id:                 q.Id.String(),
		Type:               q.VoucherType,
		Number:             q.Number,
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
		TransactionTime:    q.TransactionTime,
		CreatedAt:          q.CreatedAt,
		UpdatedAt:          q.UpdatedAt,
	}
}
