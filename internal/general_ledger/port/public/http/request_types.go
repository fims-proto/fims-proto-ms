package http

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AssignAuxiliaryCategoriesToAccountRequest struct {
	CategoryKeys []string `json:"categoryKeys"`
}

type CreateAuxiliaryCategoryRequest struct {
	Key   string `json:"key"`
	Title string `json:"title"`
}

type CreateAuxiliaryAccountRequest struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"description"`
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
