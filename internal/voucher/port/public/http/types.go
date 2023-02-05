package http

import (
	"time"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/command"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"

	"github.com/shopspring/decimal"
)

type Error struct {
	Message string `json:"message,omitempty"`
	Slug    string `json:"slug,omitempty"`
}

type CreateVoucherRequest struct {
	HeaderText         string            `json:"headerText"`
	AttachmentQuantity int               `json:"attachmentQuantity"`
	Creator            string            `json:"creator"`
	VoucherType        string            `json:"voucherType"`
	TransactionTime    time.Time         `json:"transactionTime"`
	LineItems          []LineItemRequest `json:"lineItems"`
}

type LineItemRequest struct {
	Id            uuid.UUID       `json:"id"`
	AccountNumber string          `json:"accountNumber"`
	Text          string          `json:"text"`
	Credit        decimal.Decimal `json:"credit"`
	Debit         decimal.Decimal `json:"debit"`
}

type AuditVoucherRequest struct {
	Auditor uuid.UUID `json:"auditor"`
}

type ReviewVoucherRequest struct {
	Reviewer uuid.UUID `json:"reviewer"`
}

type PostVoucherRequest struct {
	Poster uuid.UUID `json:"poster"`
}

type UpdateVoucherRequest struct {
	HeaderText      string            `json:"headerText"`
	TransactionTime time.Time         `json:"transactionTime"`
	LineItems       []LineItemRequest `json:"lineItems"`
	Updater         uuid.UUID         `json:"updater"`
}

type LineItemResponse struct {
	Id            uuid.UUID       `json:"id"`
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

type VoucherResponse struct {
	SobId              uuid.UUID          `json:"sobId"`
	Id                 uuid.UUID          `json:"id"`
	Period             PeriodResponse     `json:"period"`
	HeaderText         string             `json:"headerText"`
	DocumentNumber     string             `json:"documentNumber"`
	VoucherType        string             `json:"voucherType"`
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
		Id:            r.Id,
		Text:          r.Text,
		AccountNumber: r.AccountNumber,
		Debit:         r.Debit,
		Credit:        r.Credit,
	}
}

func (r CreateVoucherRequest) mapToCommand(sobId uuid.UUID) command.CreateVoucherCmd {
	var itemCmd []command.LineItemCmd
	for _, item := range r.LineItems {
		itemCmd = append(itemCmd, item.mapToCommand())
	}
	return command.CreateVoucherCmd{
		VoucherId:          uuid.New(),
		SobId:              sobId,
		HeaderText:         r.HeaderText,
		VoucherType:        r.VoucherType,
		AttachmentQuantity: r.AttachmentQuantity,
		LineItems:          itemCmd,
		Creator:            uuid.MustParse(r.Creator),
		TransactionTime:    r.TransactionTime,
	}
}

func lineItemDTOToVO(q query.LineItem) LineItemResponse {
	return LineItemResponse{
		Id:            q.Id,
		AccountId:     q.AccountId,
		AccountNumber: q.AccountNumber,
		Text:          q.Text,
		Debit:         q.Debit,
		Credit:        q.Credit,
		CreatedAt:     q.CreatedAt,
		UpdatedAt:     q.UpdatedAt,
	}
}

func VoucherDTOToVO(q query.Voucher) VoucherResponse {
	var itemRes []LineItemResponse
	for _, item := range q.LineItems {
		itemRes = append(itemRes, lineItemDTOToVO(item))
	}
	return VoucherResponse{
		SobId: q.SobId,
		Id:    q.Id,
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
		VoucherType:        q.VoucherType,
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
