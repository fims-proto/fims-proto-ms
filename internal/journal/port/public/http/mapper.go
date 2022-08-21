package http

import (
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/journal/app/command"
	"github/fims-proto/fims-proto-ms/internal/journal/app/query"
)

func (r LineItemRequest) mapToCommand() command.LineItemCmd {
	return command.LineItemCmd{
		ItemId:        uuid.MustParse(r.ItemId),
		Text:          r.Text,
		AccountNumber: r.AccountNumber,
		Debit:         r.Debit,
		Credit:        r.Credit,
	}
}

func (r CreateJournalEntryRequest) mapToCommand(sobId uuid.UUID) command.CreateJournalEntryCmd {
	var itemCmd []command.LineItemCmd
	for _, item := range r.LineItems {
		itemCmd = append(itemCmd, item.mapToCommand())
	}
	return command.CreateJournalEntryCmd{
		EntryId:            uuid.New(),
		SobId:              sobId,
		HeaderText:         r.HeaderText,
		JournalType:        r.JournalType,
		AttachmentQuantity: r.AttachmentQuantity,
		LineItems:          itemCmd,
		Creator:            uuid.MustParse(r.Creator),
		TransactionTime:    r.TransactionTime,
	}
}

func lineItemDTOToVO(q query.LineItem) LineItemResponse {
	return LineItemResponse{
		ItemId:        q.ItemId.String(),
		AccountId:     q.AccountId.String(),
		AccountNumber: q.AccountNumber,
		Text:          q.Text,
		Debit:         q.Debit,
		Credit:        q.Credit,
		CreatedAt:     q.CreatedAt,
		UpdatedAt:     q.UpdatedAt,
	}
}

func JournalEntryDTOToVO(q query.JournalEntry) JournalEntryResponse {
	var itemRes []LineItemResponse
	for _, item := range q.LineItems {
		itemRes = append(itemRes, lineItemDTOToVO(item))
	}
	return JournalEntryResponse{
		SobId:   q.SobId.String(),
		EntryId: q.EntryId.String(),
		Period: PeriodResponse{
			Id:            q.Period.PeriodId.String(),
			FinancialYear: q.Period.FinancialYear,
			Number:        q.Period.Number,
			OpeningTime:   q.Period.OpeningTime,
			EndingTime:    q.Period.EndingTime,
			IsClosed:      q.Period.IsClosed,
			CreatedAt:     q.Period.CreatedAt,
			UpdatedAt:     q.Period.UpdatedAt,
		},
		JournalType:        q.JournalType,
		DocumentNumber:     q.DocumentNumber,
		AttachmentQuantity: q.AttachmentQuantity,
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
		Poster: UserResponse{
			Id:     q.Poster.Id.String(),
			Traits: q.Poster.Traits,
		},
		IsReviewed:      q.IsReviewed,
		IsAudited:       q.IsAudited,
		IsPosted:        q.IsPosted,
		TransactionTime: q.TransactionTime,
		LineItems:       itemRes,
		CreatedAt:       q.CreatedAt,
		UpdatedAt:       q.UpdatedAt,
	}
}
