package http

import (
	"time"

	"github.com/shopspring/decimal"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
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
	ItemId        string          `json:"itemId"`
	AccountNumber string          `json:"accountNumber"`
	Text          string          `json:"text"`
	Credit        decimal.Decimal `json:"credit"`
	Debit         decimal.Decimal `json:"debit"`
}

type AuditJournalEntryRequest struct {
	Auditor string `json:"auditor"`
}

type ReviewJournalEntryRequest struct {
	Reviewer string `json:"reviewer"`
}

type PostJournalEntryRequest struct {
	Poster string `json:"poster"`
}

type UpdateJournalEntryRequest struct {
	TransactionTime time.Time         `json:"transactionTime"`
	LineItems       []LineItemRequest `json:"lineItems"`
}

type LineItemResponse struct {
	ItemId        string          `json:"itemId"`
	AccountId     string          `json:"accountId"`
	AccountNumber string          `json:"accountNumber"`
	Text          string          `json:"text"`
	Credit        decimal.Decimal `json:"credit"`
	Debit         decimal.Decimal `json:"debit"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

type UserResponse struct {
	Id     string `json:"id"`
	Traits any    `json:"traits"`
}

type PeriodResponse struct {
	Id            string    `json:"id"`
	FinancialYear int       `json:"financialYear"`
	Number        int       `json:"number"`
	OpeningTime   time.Time `json:"openingTime"`
	EndingTime    time.Time `json:"endingTime"`
	IsClosed      bool      `json:"isClosed"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type JournalEntryResponse struct {
	SobId              string             `json:"sobId"`
	EntryId            string             `json:"entryId"`
	Period             PeriodResponse     `json:"period"`
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
