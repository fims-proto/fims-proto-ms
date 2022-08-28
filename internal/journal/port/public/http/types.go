package http

import (
	"time"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/journal/app/command"
	"github/fims-proto/fims-proto-ms/internal/journal/app/query"

	"github.com/shopspring/decimal"
)

type Error struct {
	Message string `json:"message,omitempty"`
	Slug    string `json:"slug,omitempty"`
}

type CreateJournalEntryRequest struct {
	HeaderText         string            `json:"headerText"`
	AttachmentQuantity int               `json:"attachmentQuantity"`
	Creator            string            `json:"creator"`
	JournalType        string            `json:"journalType"`
	TransactionTime    time.Time         `json:"transactionTime"`
	LineItems          []LineItemRequest `json:"lineItems"`
}

type LineItemRequest struct {
	ItemId        uuid.UUID       `json:"itemId"`
	AccountNumber string          `json:"accountNumber"`
	Text          string          `json:"text"`
	Credit        decimal.Decimal `json:"credit"`
	Debit         decimal.Decimal `json:"debit"`
}

type AuditJournalEntryRequest struct {
	Auditor uuid.UUID `json:"auditor"`
}

type ReviewJournalEntryRequest struct {
	Reviewer uuid.UUID `json:"reviewer"`
}

type PostJournalEntryRequest struct {
	Poster uuid.UUID `json:"poster"`
}

type UpdateJournalEntryRequest struct {
	TransactionTime time.Time         `json:"transactionTime"`
	LineItems       []LineItemRequest `json:"lineItems"`
}

type LineItemResponse struct {
	ItemId        uuid.UUID       `json:"itemId"`
	AccountId     uuid.UUID       `json:"accountId"`
	AccountNumber string          `json:"accountNumber"`
	Text          string          `json:"text"`
	Credit        decimal.Decimal `json:"credit"`
	Debit         decimal.Decimal `json:"debit"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

type UserResponse struct {
	Id     uuid.UUID `json:"id"`
	Traits any       `json:"traits"`
}

type PeriodResponse struct {
	Id            uuid.UUID `json:"id"`
	FinancialYear int       `json:"financialYear"`
	Number        int       `json:"number"`
	OpeningTime   time.Time `json:"openingTime"`
	EndingTime    time.Time `json:"endingTime"`
	IsClosed      bool      `json:"isClosed"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type JournalEntryResponse struct {
	SobId              uuid.UUID          `json:"sobId"`
	EntryId            uuid.UUID          `json:"entryId"`
	Period             PeriodResponse     `json:"period"`
	HeaderText         string             `json:"headerText"`
	DocumentNumber     string             `json:"documentNumber"`
	JournalType        string             `json:"journalType"`
	AttachmentQuantity int                `json:"attachmentQuantity"`
	Creator            UserResponse       `json:"creator"`
	Auditor            UserResponse       `json:"auditor"`
	Reviewer           UserResponse       `json:"reviewer"`
	Poster             UserResponse       `json:"poster"`
	Credit             decimal.Decimal    `json:"credit"`
	Debit              decimal.Decimal    `json:"debit"`
	IsAudited          bool               `json:"isAudited"`
	IsPosted           bool               `json:"isPosted"`
	IsReviewed         bool               `json:"isReviewed"`
	TransactionTime    time.Time          `json:"transactionTime"`
	LineItems          []LineItemResponse `json:"lineItems"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

// mapper

func (r LineItemRequest) mapToCommand() command.LineItemCmd {
	return command.LineItemCmd{
		ItemId:        r.ItemId,
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
		ItemId:        q.ItemId,
		AccountId:     q.AccountId,
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
		SobId:   q.SobId,
		EntryId: q.EntryId,
		Period: PeriodResponse{
			Id:            q.Period.PeriodId,
			FinancialYear: q.Period.FinancialYear,
			Number:        q.Period.Number,
			OpeningTime:   q.Period.OpeningTime,
			EndingTime:    q.Period.EndingTime,
			IsClosed:      q.Period.IsClosed,
			CreatedAt:     q.Period.CreatedAt,
			UpdatedAt:     q.Period.UpdatedAt,
		},
		HeaderText:         q.HeaderText,
		JournalType:        q.JournalType,
		DocumentNumber:     q.DocumentNumber,
		AttachmentQuantity: q.AttachmentQuantity,
		Debit:              q.Debit,
		Credit:             q.Credit,
		Creator: UserResponse{
			Id:     q.Creator.Id,
			Traits: q.Creator.Traits,
		},
		Reviewer: UserResponse{
			Id:     q.Reviewer.Id,
			Traits: q.Reviewer.Traits,
		},
		Auditor: UserResponse{
			Id:     q.Auditor.Id,
			Traits: q.Auditor.Traits,
		},
		Poster: UserResponse{
			Id:     q.Poster.Id,
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
