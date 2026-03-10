package http

import (
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateAccountRequest struct {
	Title                 string   `json:"title"`
	LevelNumber           int      `json:"levelNumber"`
	SuperiorAccountNumber string   `json:"superiorAccountNumber,omitempty"`
	BalanceDirection      string   `json:"balanceDirection"`
	Class                 string   `json:"class,omitempty"`
	Group                 string   `json:"group,omitempty"`
	CategoryKeys          []string `json:"categoryKeys,omitempty"`
}

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

type CreateJournalRequest struct {
	HeaderText         string                           `json:"headerText"`
	AttachmentQuantity int                              `json:"attachmentQuantity"`
	Creator            string                           `json:"creator"`
	JournalType        string                           `json:"journalType"`
	TransactionDate    transaction_date.TransactionDate `json:"transactionDate"`
	JournalLines       []JournalLineRequest             `json:"journalLines"`
}

type JournalLineRequest struct {
	Id                uuid.UUID              `json:"id"`
	AccountNumber     string                 `json:"accountNumber"`
	AuxiliaryAccounts []AuxiliaryItemRequest `json:"auxiliaryAccounts"`
	Text              string                 `json:"text"`
	Amount            decimal.Decimal        `json:"amount"`
}

type AuxiliaryItemRequest struct {
	CategoryKey string `json:"categoryKey"`
	AccountKey  string `json:"accountKey"`
}

type AuditJournalRequest struct {
	Auditor uuid.UUID `json:"auditor"`
}

type ReviewJournalRequest struct {
	Reviewer uuid.UUID `json:"reviewer"`
}

type PostJournalRequest struct {
	Poster uuid.UUID `json:"poster"`
}

type UpdateJournalRequest struct {
	HeaderText      string                           `json:"headerText"`
	TransactionDate transaction_date.TransactionDate `json:"transactionDate"`
	JournalLines    []JournalLineRequest             `json:"journalLines"`
	Updater         uuid.UUID                        `json:"updater"`
}

type InitializeLedgersBalanceRequest struct {
	Ledgers []InitializeLedgersBalanceItemRequest `json:"ledgers" binding:"required"`
}

type InitializeLedgersBalanceItemRequest struct {
	AccountNumber  string          `json:"accountNumber"`
	OpeningBalance decimal.Decimal `json:"openingBalance"`
}

// mapper

func (r JournalLineRequest) mapToCommand() command.JournalLineCmd {
	var auxiliaryItemCmds []command.AuxiliaryItemCmd
	for _, auxiliaryAccount := range r.AuxiliaryAccounts {
		auxiliaryItemCmds = append(auxiliaryItemCmds, command.AuxiliaryItemCmd{
			CategoryKey: auxiliaryAccount.CategoryKey,
			AccountKey:  auxiliaryAccount.AccountKey,
		})
	}

	return command.JournalLineCmd{
		Id:                r.Id,
		Text:              r.Text,
		AccountNumber:     r.AccountNumber,
		AuxiliaryAccounts: auxiliaryItemCmds,
		Amount:            r.Amount,
	}
}

func (r CreateJournalRequest) mapToCommand(sobId uuid.UUID) command.CreateJournalCmd {
	var itemCmd []command.JournalLineCmd
	for _, item := range r.JournalLines {
		itemCmd = append(itemCmd, item.mapToCommand())
	}
	return command.CreateJournalCmd{
		JournalId:          uuid.New(),
		SobId:              sobId,
		HeaderText:         r.HeaderText,
		JournalType:        r.JournalType,
		AttachmentQuantity: r.AttachmentQuantity,
		JournalLines:       itemCmd,
		Creator:            uuid.MustParse(r.Creator),
		TransactionDate:    r.TransactionDate,
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
