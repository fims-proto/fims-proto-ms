package http

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UpdateAccountRequest struct {
	Title            string   `json:"title,omitempty"`
	LevelNumber      int      `json:"levelNumber,omitempty"`
	BalanceDirection string   `json:"balanceDirection,omitempty"`
	Group            string   `json:"group"`
	CategoryKeys     []string `json:"categoryKeys,omitempty"`
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
	Id                uuid.UUID              `json:"id"`
	AccountNumber     string                 `json:"accountNumber"`
	AuxiliaryAccounts []AuxiliaryItemRequest `json:"auxiliaryAccounts"`
	Text              string                 `json:"text"`
	Credit            decimal.Decimal        `json:"credit"`
	Debit             decimal.Decimal        `json:"debit"`
}

type AuxiliaryItemRequest struct {
	CategoryKey string `json:"categoryKey"`
	AccountKey  string `json:"accountKey"`
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

type InitializeLedgersBalanceRequest struct {
	Ledgers []InitializeLedgersBalanceItemRequest `json:"ledgers" binding:"required"`
}

type InitializeLedgersBalanceItemRequest struct {
	AccountNumber  string          `json:"accountNumber"`
	OpeningBalance decimal.Decimal `json:"openingBalance"`
}

// mapper

func (r LineItemRequest) mapToCommand() command.LineItemCmd {
	var auxiliaryItemCmds []command.AuxiliaryItemCmd
	for _, auxiliaryAccount := range r.AuxiliaryAccounts {
		auxiliaryItemCmds = append(auxiliaryItemCmds, command.AuxiliaryItemCmd{
			CategoryKey: auxiliaryAccount.CategoryKey,
			AccountKey:  auxiliaryAccount.AccountKey,
		})
	}

	return command.LineItemCmd{
		Id:                r.Id,
		Text:              r.Text,
		AccountNumber:     r.AccountNumber,
		AuxiliaryAccounts: auxiliaryItemCmds,
		Debit:             r.Debit,
		Credit:            r.Credit,
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

func (r InitializeLedgersBalanceRequest) mapToCommand(sobId uuid.UUID) command.InitializeLedgersBalanceCmd {
	var itemCmd []command.InitializeLedgersBalanceItemCmd
	for _, l := range r.Ledgers {
		itemCmd = append(itemCmd, command.InitializeLedgersBalanceItemCmd{
			AccountNumber:  l.AccountNumber,
			OpeningBalance: l.OpeningBalance,
		})
	}

	return command.InitializeLedgersBalanceCmd{
		SobId:   sobId,
		Ledgers: itemCmd,
	}
}
