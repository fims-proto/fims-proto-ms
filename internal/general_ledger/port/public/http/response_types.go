package http

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type AccountResponse struct {
	Id                  uuid.UUID                   `json:"id,omitempty"`
	SobId               uuid.UUID                   `json:"sobId,omitempty"`
	SuperiorAccountId   uuid.UUID                   `json:"superiorAccountId,omitempty"`
	Title               string                      `json:"title,omitempty"`
	AccountNumber       string                      `json:"accountNumber,omitempty"`
	NumberHierarchy     []int                       `json:"numberHierarchy,omitempty"`
	Level               int                         `json:"level,omitempty"`
	AccountType         string                      `json:"accountType,omitempty"`
	BalanceDirection    string                      `json:"balanceDirection,omitempty"`
	AuxiliaryCategories []AuxiliaryCategoryResponse `json:"auxiliaryCategories"`
	CreatedAt           time.Time                   `json:"createdAt"`
	UpdatedAt           time.Time                   `json:"updatedAt"`
}

type AuxiliaryCategoryResponse struct {
	Id         uuid.UUID `json:"id,omitempty"`
	SobId      uuid.UUID `json:"sobId,omitempty"`
	Key        string    `json:"key,omitempty"`
	Title      string    `json:"title,omitempty"`
	IsStandard bool      `json:"isStandard"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type AuxiliaryAccountResponse struct {
	Id          uuid.UUID                 `json:"id,omitempty"`
	Category    AuxiliaryCategoryResponse `json:"category"`
	Key         string                    `json:"key,omitempty"`
	Title       string                    `json:"title,omitempty"`
	Description string                    `json:"description,omitempty"`
	CreatedAt   time.Time                 `json:"createdAt"`
	UpdatedAt   time.Time                 `json:"updatedAt"`
}

type PeriodResponse struct {
	Id           uuid.UUID `json:"id,omitempty"`
	SobId        uuid.UUID `json:"sobId,omitempty"`
	FiscalYear   int       `json:"fiscalYear,omitempty"`
	PeriodNumber int       `json:"periodNumber,omitempty"`
	OpeningTime  time.Time `json:"openingTime"`
	EndingTime   time.Time `json:"endingTime"`
	IsClosed     bool      `json:"isClosed"`
	IsCurrent    bool      `json:"isCurrent"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type LedgerResponse struct {
	Id             uuid.UUID       `json:"id,omitempty"`
	SobId          uuid.UUID       `json:"sobId,omitempty"`
	AccountId      uuid.UUID       `json:"accountId,omitempty"`
	PeriodId       uuid.UUID       `json:"periodId,omitempty"`
	OpeningBalance decimal.Decimal `json:"openingBalance"`
	EndingBalance  decimal.Decimal `json:"endingBalance"`
	PeriodDebit    decimal.Decimal `json:"periodDebit"`
	PeriodCredit   decimal.Decimal `json:"periodCredit"`
	Account        AccountResponse `json:"account"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type AuxiliaryLedgerResponse struct {
	Id               uuid.UUID                `json:"id,omitempty"`
	PeriodId         uuid.UUID                `json:"periodId,omitempty"`
	AuxiliaryAccount AuxiliaryAccountResponse `json:"auxiliaryAccount"`
	OpeningBalance   decimal.Decimal          `json:"openingBalance"`
	EndingBalance    decimal.Decimal          `json:"endingBalance"`
	PeriodDebit      decimal.Decimal          `json:"periodDebit"`
	PeriodCredit     decimal.Decimal          `json:"periodCredit"`
	CreatedAt        time.Time                `json:"createdAt"`
	UpdatedAt        time.Time                `json:"updatedAt"`
}

type LineItemResponse struct {
	Id                uuid.UUID                  `json:"id,omitempty"`
	AccountId         uuid.UUID                  `json:"accountId,omitempty"`
	AccountNumber     string                     `json:"accountNumber,omitempty"`
	AuxiliaryAccounts []AuxiliaryAccountResponse `json:"auxiliaryAccounts,omitempty"`
	Text              string                     `json:"text,omitempty"`
	Credit            decimal.Decimal            `json:"credit"`
	Debit             decimal.Decimal            `json:"debit"`
	CreatedAt         time.Time                  `json:"createdAt"`
	UpdatedAt         time.Time                  `json:"updatedAt"`
}

type VoucherResponse struct {
	Id                 uuid.UUID          `json:"id,omitempty"`
	SobId              uuid.UUID          `json:"sobId,omitempty"`
	Period             PeriodResponse     `json:"period"`
	HeaderText         string             `json:"headerText,omitempty"`
	DocumentNumber     string             `json:"documentNumber,omitempty"`
	VoucherType        string             `json:"voucherType,omitempty"`
	AttachmentQuantity int                `json:"attachmentQuantity,omitempty"`
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
	LineItems          []LineItemResponse `json:"lineItems,omitempty"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

type UserResponse struct {
	Id     uuid.UUID `json:"id"`
	Traits any       `json:"traits"`
}

// mapper

func accountDTOToVO(dto query.Account) AccountResponse {
	var categories []AuxiliaryCategoryResponse
	for _, category := range dto.AuxiliaryCategories {
		categories = append(categories, auxiliaryCategoryDTOToVO(category))
	}
	return AccountResponse{
		Id:                  dto.Id,
		SobId:               dto.SobId,
		SuperiorAccountId:   dto.SuperiorAccountId,
		Title:               dto.Title,
		AccountNumber:       dto.AccountNumber,
		NumberHierarchy:     dto.NumberHierarchy,
		Level:               dto.Level,
		AccountType:         dto.AccountType,
		BalanceDirection:    dto.BalanceDirection,
		AuxiliaryCategories: categories,
		CreatedAt:           dto.CreatedAt,
		UpdatedAt:           dto.UpdatedAt,
	}
}

func auxiliaryCategoryDTOToVO(dto query.AuxiliaryCategory) AuxiliaryCategoryResponse {
	return AuxiliaryCategoryResponse(dto)
}

func auxiliaryAccountDTOToVO(dto query.AuxiliaryAccount) AuxiliaryAccountResponse {
	return AuxiliaryAccountResponse{
		Id:          dto.Id,
		Category:    auxiliaryCategoryDTOToVO(dto.Category),
		Key:         dto.Key,
		Title:       dto.Title,
		Description: dto.Description,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
	}
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

func auxiliaryLedgerDTOToVO(dto query.AuxiliaryLedger) AuxiliaryLedgerResponse {
	return AuxiliaryLedgerResponse{
		Id:               dto.Id,
		PeriodId:         dto.PeriodId,
		AuxiliaryAccount: auxiliaryAccountDTOToVO(dto.AuxiliaryAccount),
		OpeningBalance:   dto.OpeningBalance,
		EndingBalance:    dto.EndingBalance,
		PeriodDebit:      dto.PeriodDebit,
		PeriodCredit:     dto.PeriodCredit,
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
	}
}

func lineItemDTOToVO(dto query.LineItem) LineItemResponse {
	var auxiliaryAccounts []AuxiliaryAccountResponse
	for _, auxiliaryAccount := range dto.AuxiliaryAccounts {
		auxiliaryAccounts = append(auxiliaryAccounts, auxiliaryAccountDTOToVO(auxiliaryAccount))
	}
	return LineItemResponse{
		Id:                dto.Id,
		AccountId:         dto.AccountId,
		AccountNumber:     dto.AccountNumber,
		AuxiliaryAccounts: auxiliaryAccounts,
		Text:              dto.Text,
		Credit:            dto.Credit,
		Debit:             dto.Debit,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
	}
}

func VoucherDTOToVO(dto query.Voucher) VoucherResponse {
	var itemRes []LineItemResponse
	for _, item := range dto.LineItems {
		itemRes = append(itemRes, lineItemDTOToVO(item))
	}
	return VoucherResponse{
		SobId:              dto.SobId,
		Id:                 dto.Id,
		Period:             periodDTOToVO(dto.Period),
		HeaderText:         dto.HeaderText,
		VoucherType:        dto.VoucherType,
		DocumentNumber:     dto.DocumentNumber,
		AttachmentQuantity: dto.AttachmentQuantity,
		Debit:              dto.Debit,
		Credit:             dto.Credit,
		Creator: UserResponse{
			Id:     dto.Creator.Id,
			Traits: dto.Creator.Traits,
		},
		Reviewer: UserResponse{
			Id:     dto.Reviewer.Id,
			Traits: dto.Reviewer.Traits,
		},
		Auditor: UserResponse{
			Id:     dto.Auditor.Id,
			Traits: dto.Auditor.Traits,
		},
		Poster: UserResponse{
			Id:     dto.Poster.Id,
			Traits: dto.Poster.Traits,
		},
		IsReviewed:      dto.IsReviewed,
		IsAudited:       dto.IsAudited,
		IsPosted:        dto.IsPosted,
		TransactionTime: dto.TransactionTime,
		LineItems:       itemRes,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}
