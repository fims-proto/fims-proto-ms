package http

import (
	"strconv"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type AccountClass struct {
	Class  string   `json:"id"`
	Groups []string `json:"groups"`
}

type AccountResponse struct {
	Id                  uuid.UUID                   `json:"id,omitempty"`
	SobId               uuid.UUID                   `json:"sobId,omitempty"`
	SuperiorAccountId   *uuid.UUID                  `json:"superiorAccountId,omitempty"`
	Title               string                      `json:"title,omitempty"`
	AccountNumber       string                      `json:"accountNumber,omitempty"`
	NumberHierarchy     []int                       `json:"numberHierarchy,omitempty"`
	Level               int                         `json:"level"`
	IsLeaf              bool                        `json:"isLeaf"`
	Class               string                      `json:"class"`
	Group               string                      `json:"group"`
	BalanceDirection    string                      `json:"balanceDirection,omitempty"`
	AuxiliaryCategories []AuxiliaryCategoryResponse `json:"auxiliaryCategories,omitempty"`
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
	FiscalYear   int       `json:"fiscalYear"`
	PeriodNumber int       `json:"periodNumber"`
	IsClosed     bool      `json:"isClosed"`
	IsCurrent    bool      `json:"isCurrent"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type LedgerResponse struct {
	Id            uuid.UUID       `json:"id,omitempty"`
	SobId         uuid.UUID       `json:"sobId,omitempty"`
	AccountId     uuid.UUID       `json:"accountId,omitempty"`
	PeriodId      uuid.UUID       `json:"periodId,omitempty"`
	OpeningAmount decimal.Decimal `json:"openingAmount"`
	PeriodAmount  decimal.Decimal `json:"periodAmount"`
	PeriodDebit   decimal.Decimal `json:"periodDebit"`
	PeriodCredit  decimal.Decimal `json:"periodCredit"`
	EndingAmount  decimal.Decimal `json:"endingAmount"`
	Account       AccountResponse `json:"account"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

type LedgerSummaryResponse struct {
	AccountId     uuid.UUID       `json:"accountId"`
	PeriodId      uuid.UUID       `json:"periodId"`
	OpeningAmount decimal.Decimal `json:"openingAmount"`
	PeriodAmount  decimal.Decimal `json:"periodAmount"`
	PeriodDebit   decimal.Decimal `json:"periodDebit"`
	PeriodCredit  decimal.Decimal `json:"periodCredit"`
	EndingAmount  decimal.Decimal `json:"endingAmount"`
}

type AuxiliaryLedgerSummaryResponse struct {
	AuxiliaryAccountId    uuid.UUID       `json:"auxiliaryAccountId"`
	AuxiliaryAccountTitle string          `json:"auxiliaryAccountTitle"`
	OpeningAmount         decimal.Decimal `json:"openingAmount"`
	PeriodAmount          decimal.Decimal `json:"periodAmount"`
	PeriodDebit           decimal.Decimal `json:"periodDebit"`
	PeriodCredit          decimal.Decimal `json:"periodCredit"`
	EndingAmount          decimal.Decimal `json:"endingAmount"`
}

type PeriodAndLedgersResponse struct {
	Period  PeriodResponse   `json:"period"`
	Ledgers []LedgerResponse `json:"ledgers"`
}

type AuxiliaryLedgerResponse struct {
	Id                   uuid.UUID                 `json:"id,omitempty"`
	SobId                uuid.UUID                 `json:"sobId,omitempty"`
	PeriodId             uuid.UUID                 `json:"periodId,omitempty"`
	Account              AccountResponse           `json:"account"`
	AuxiliaryCategory    AuxiliaryCategoryResponse `json:"auxiliaryCategory"`
	AuxiliaryAccount     AuxiliaryAccountResponse  `json:"auxiliaryAccount"`
	OpeningDebitBalance  decimal.Decimal           `json:"openingBalance"`
	OpeningCreditBalance decimal.Decimal           `json:"openingCreditBalance"`
	PeriodDebit          decimal.Decimal           `json:"periodDebit"`
	PeriodCredit         decimal.Decimal           `json:"periodCredit"`
	EndingDebitBalance   decimal.Decimal           `json:"endingBalance"`
	EndingCreditBalance  decimal.Decimal           `json:"endingCreditBalance"`
	CreatedAt            time.Time                 `json:"createdAt"`
	UpdatedAt            time.Time                 `json:"updatedAt"`
}

type LineItemResponse struct {
	Id                uuid.UUID                  `json:"id,omitempty"`
	Account           AccountResponse            `json:"account"`
	AuxiliaryAccounts []AuxiliaryAccountResponse `json:"auxiliaryAccounts,omitempty"`
	Text              string                     `json:"text,omitempty"`
	Amount            decimal.Decimal            `json:"amount"`
	CreatedAt         time.Time                  `json:"createdAt"`
	UpdatedAt         time.Time                  `json:"updatedAt"`
}

type VoucherResponse struct {
	Id                 uuid.UUID                        `json:"id,omitempty"`
	SobId              uuid.UUID                        `json:"sobId,omitempty"`
	Period             PeriodResponse                   `json:"period"`
	HeaderText         string                           `json:"headerText,omitempty"`
	DocumentNumber     string                           `json:"documentNumber,omitempty"`
	VoucherType        string                           `json:"voucherType,omitempty"`
	AttachmentQuantity int                              `json:"attachmentQuantity"`
	Creator            *UserResponse                    `json:"creator"`
	Auditor            *UserResponse                    `json:"auditor"`
	Reviewer           *UserResponse                    `json:"reviewer"`
	Poster             *UserResponse                    `json:"poster"`
	Amount             decimal.Decimal                  `json:"amount"`
	IsAudited          bool                             `json:"isAudited"`
	IsPosted           bool                             `json:"isPosted"`
	IsReviewed         bool                             `json:"isReviewed"`
	TransactionDate    transaction_date.TransactionDate `json:"transactionDate" swaggertype:"string"`
	LineItems          []LineItemResponse               `json:"lineItems,omitempty"`
	CreatedAt          time.Time                        `json:"createdAt"`
	UpdatedAt          time.Time                        `json:"updatedAt"`
}

type UserResponse struct {
	Id     uuid.UUID `json:"id"`
	Traits any       `json:"traits"`
}

type LedgerEntryResponse struct {
	VoucherId       uuid.UUID                        `json:"voucherId"`
	VoucherNumber   string                           `json:"voucherNumber"`
	TransactionDate transaction_date.TransactionDate `json:"transactionDate" swaggertype:"string"`
	Text            string                           `json:"text"`
	Amount          decimal.Decimal                  `json:"amount"`
	CreatedAt       time.Time                        `json:"createdAt"`
	UpdatedAt       time.Time                        `json:"updatedAt"`
}

// mapper

func accountDTOToVO(dto query.Account) AccountResponse {
	return AccountResponse{
		Id:                  dto.Id,
		SobId:               dto.SobId,
		SuperiorAccountId:   dto.SuperiorAccountId,
		Title:               dto.Title,
		AccountNumber:       dto.AccountNumber,
		NumberHierarchy:     dto.NumberHierarchy,
		Level:               dto.Level,
		IsLeaf:              dto.IsLeaf,
		Class:               strconv.Itoa(dto.Class),
		Group:               strconv.Itoa(dto.Group),
		BalanceDirection:    dto.BalanceDirection,
		AuxiliaryCategories: converter.DTOsToVOs(dto.AuxiliaryCategories, auxiliaryCategoryDTOToVO),
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
		Id:            dto.Id,
		SobId:         dto.SobId,
		AccountId:     dto.AccountId,
		PeriodId:      dto.PeriodId,
		OpeningAmount: dto.OpeningAmount,
		PeriodAmount:  dto.PeriodAmount,
		PeriodDebit:   dto.PeriodDebit,
		PeriodCredit:  dto.PeriodCredit,
		EndingAmount:  dto.EndingAmount,
		Account:       accountDTOToVO(dto.Account),
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
	}
}

func ledgerSummaryToVO(dto query.LedgerSummary) LedgerSummaryResponse {
	return LedgerSummaryResponse(dto)
}

func auxiliaryLedgerSummaryToVO(dto query.AuxiliaryLedgerSummary) AuxiliaryLedgerSummaryResponse {
	return AuxiliaryLedgerSummaryResponse(dto)
}

func lineItemDTOToVO(dto query.LineItem) LineItemResponse {
	return LineItemResponse{
		Id:                dto.Id,
		Account:           accountDTOToVO(dto.Account),
		AuxiliaryAccounts: converter.DTOsToVOs(dto.AuxiliaryAccounts, auxiliaryAccountDTOToVO),
		Text:              dto.Text,
		Amount:            dto.Amount,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
	}
}

func voucherDTOToVO(dto query.Voucher) VoucherResponse {
	userOrNil := func(u *query.User) *UserResponse {
		if u != nil {
			return &UserResponse{
				Id:     u.Id,
				Traits: u.Traits,
			}
		}
		return nil
	}

	return VoucherResponse{
		SobId:              dto.SobId,
		Id:                 dto.Id,
		Period:             periodDTOToVO(dto.Period),
		HeaderText:         dto.HeaderText,
		VoucherType:        dto.VoucherType,
		DocumentNumber:     dto.DocumentNumber,
		AttachmentQuantity: dto.AttachmentQuantity,
		Amount:             dto.Amount,
		Creator:            userOrNil(dto.Creator),
		Reviewer:           userOrNil(dto.Reviewer),
		Auditor:            userOrNil(dto.Auditor),
		Poster:             userOrNil(dto.Poster),
		IsReviewed:         dto.IsReviewed,
		IsAudited:          dto.IsAudited,
		IsPosted:           dto.IsPosted,
		TransactionDate:    dto.TransactionDate,
		LineItems:          converter.DTOsToVOs(dto.LineItems, lineItemDTOToVO),
		CreatedAt:          dto.CreatedAt,
		UpdatedAt:          dto.UpdatedAt,
	}
}

func ledgerEntryDTOToVO(dto query.LedgerEntry) LedgerEntryResponse {
	return LedgerEntryResponse(dto)
}
