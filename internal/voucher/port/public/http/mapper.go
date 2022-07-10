package http

import (
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
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
		Creator:            uuid.MustParse(r.Creator),
		TransactionTime:    r.TransactionTime,
	}
}

func mapFromLineItemQuery(q query.LineItem) LineItemResponse {
	return LineItemResponse{
		Id:            q.Id.String(),
		AccountId:     q.AccountId.String(),
		AccountNumber: q.AccountNumber,
		Summary:       q.Summary,
		Debit:         q.Debit,
		Credit:        q.Credit,
	}
}

func mapFromVoucherQuery(q query.Voucher) VoucherResponse {
	var itemRes []LineItemResponse
	for _, item := range q.LineItems {
		itemRes = append(itemRes, mapFromLineItemQuery(item))
	}
	return VoucherResponse{
		SobId: q.SobId.String(),
		Id:    q.Id.String(),
		Period: PeriodResponse{
			Id:            q.Period.Id.String(),
			FinancialYear: q.Period.FinancialYear,
			Number:        q.Period.Number,
			OpeningTime:   q.Period.OpeningTime,
			EndingTime:    q.Period.EndingTime,
			IsClosed:      q.Period.IsClosed,
		},
		Type:               q.VoucherType,
		Number:             q.Number,
		AttachmentQuantity: int(q.AttachmentQuantity),
		LineItems:          itemRes,
		Debit:              q.Debit,
		Credit:             q.Credit,
		Creator: UserResponse{
			Id:     q.Creator.Id.String(),
			Traits: q.Creator.Traits,
		},
		Reviewer: UserResponse{
			Id:     q.Reviewer.Id.String(),
			Traits: q.Reviewer.Traits,
		},
		Auditor: UserResponse{
			Id:     q.Auditor.Id.String(),
			Traits: q.Auditor.Traits,
		},
		IsReviewed:      q.IsReviewed,
		IsAudited:       q.IsAudited,
		IsPosted:        q.IsPosted,
		TransactionTime: q.TransactionTime,
		CreatedAt:       q.CreatedAt,
		UpdatedAt:       q.UpdatedAt,
	}
}
