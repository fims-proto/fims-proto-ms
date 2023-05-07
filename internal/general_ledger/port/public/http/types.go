package http

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

// requests

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

// responses

type AccountResponse struct {
	Id                uuid.UUID `json:"id"`
	SobId             uuid.UUID `json:"sobId"`
	SuperiorAccountId uuid.UUID `json:"superiorAccountId"`
	Title             string    `json:"title"`
	AccountNumber     string    `json:"accountNumber"`
	NumberHierarchy   []int     `json:"numberHierarchy"`
	Level             int       `json:"level"`
	AccountType       string    `json:"accountType"`
	BalanceDirection  string    `json:"balanceDirection"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type PeriodResponse struct {
	Id           uuid.UUID `json:"id"`
	SobId        uuid.UUID `json:"sobId"`
	FiscalYear   int       `json:"fiscalYear"`
	PeriodNumber int       `json:"periodNumber"`
	OpeningTime  time.Time `json:"openingTime"`
	EndingTime   time.Time `json:"endingTime"`
	IsClosed     bool      `json:"isClosed"`
	IsCurrent    bool      `json:"isCurrent"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type LedgerResponse struct {
	Id             uuid.UUID       `json:"id"`
	SobId          uuid.UUID       `json:"sobId"`
	AccountId      uuid.UUID       `json:"accountId"`
	PeriodId       uuid.UUID       `json:"periodId"`
	OpeningBalance decimal.Decimal `json:"openingBalance"`
	EndingBalance  decimal.Decimal `json:"endingBalance"`
	PeriodDebit    decimal.Decimal `json:"periodDebit"`
	PeriodCredit   decimal.Decimal `json:"periodCredit"`
	Account        AccountResponse `json:"account"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
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

type VoucherResponse struct {
	Id                 uuid.UUID          `json:"id"`
	SobId              uuid.UUID          `json:"sobId"`
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

func accountDTOToVO(dto query.Account) AccountResponse {
	return AccountResponse(dto)
}

func periodDTOToVO(dto query.Period) PeriodResponse {
	return PeriodResponse(dto)
}

func ledgerDTOToVO(dto query.Ledger) LedgerResponse {
	return LedgerResponse{
		Id:             dto.Id,
		SobId:          dto.SobId,
		AccountId:      dto.AccountId,
		PeriodId:       dto.PeriodId,
		OpeningBalance: dto.OpeningBalance,
		EndingBalance:  dto.EndingBalance,
		PeriodDebit:    dto.PeriodDebit,
		PeriodCredit:   dto.PeriodCredit,
		Account:        accountDTOToVO(dto.Account),
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
	}
}

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
		SobId:              q.SobId,
		Id:                 q.Id,
		Period:             periodDTOToVO(q.Period),
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
