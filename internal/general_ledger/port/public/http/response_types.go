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
	SobId             uuid.UUID       `json:"sobId,omitempty"`
	AccountId         uuid.UUID       `json:"accountId,omitempty"`
	SuperiorAccountId *uuid.UUID      `json:"superiorAccountId,omitempty"`
	AccountNumber     string          `json:"accountNumber,omitempty"`
	AccountTitle      string          `json:"accountTitle,omitempty"`
	AccountClass      string          `json:"accountClass"`
	AccountGroup      string          `json:"accountGroup"`
	BalanceDirection  string          `json:"balanceDirection,omitempty"`
	IsLeaf            bool            `json:"isLeaf"`
	OpeningAmount     decimal.Decimal `json:"openingAmount"`
	PeriodAmount      decimal.Decimal `json:"periodAmount"`
	PeriodDebit       decimal.Decimal `json:"periodDebit"`
	PeriodCredit      decimal.Decimal `json:"periodCredit"`
	EndingAmount      decimal.Decimal `json:"endingAmount"`
}

type LedgerSummaryResponse struct {
	AccountId     uuid.UUID       `json:"accountId"`
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

type JournalLineResponse struct {
	Id                uuid.UUID                  `json:"id,omitempty"`
	Account           AccountResponse            `json:"account"`
	AuxiliaryAccounts []AuxiliaryAccountResponse `json:"auxiliaryAccounts,omitempty"`
	Text              string                     `json:"text,omitempty"`
	Amount            decimal.Decimal            `json:"amount"`
	CreatedAt         time.Time                  `json:"createdAt"`
	UpdatedAt         time.Time                  `json:"updatedAt"`
}

type JournalResponse struct {
	Id                 uuid.UUID                        `json:"id,omitempty"`
	SobId              uuid.UUID                        `json:"sobId,omitempty"`
	Period             PeriodResponse                   `json:"period"`
	HeaderText         string                           `json:"headerText,omitempty"`
	DocumentNumber     string                           `json:"documentNumber,omitempty"`
	JournalType        string                           `json:"journalType,omitempty"`
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
	JournalLines       []JournalLineResponse            `json:"journalLines,omitempty"`
	CreatedAt          time.Time                        `json:"createdAt"`
	UpdatedAt          time.Time                        `json:"updatedAt"`
}

type UserResponse struct {
	Id     uuid.UUID `json:"id"`
	Traits any       `json:"traits"`
}

type LedgerEntryResponse struct {
	JournalId       uuid.UUID                        `json:"journalId"`
	JournalNumber   string                           `json:"journalNumber"`
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
		SobId:             dto.SobId,
		AccountId:         dto.AccountId,
		SuperiorAccountId: dto.Account.SuperiorAccountId,
		AccountNumber:     dto.Account.AccountNumber,
		AccountTitle:      dto.Account.Title,
		AccountClass:      strconv.Itoa(dto.Account.Class),
		AccountGroup:      strconv.Itoa(dto.Account.Group),
		BalanceDirection:  dto.Account.BalanceDirection,
		IsLeaf:            dto.Account.IsLeaf,
		OpeningAmount:     dto.OpeningAmount,
		PeriodAmount:      dto.PeriodAmount,
		PeriodDebit:       dto.PeriodDebit,
		PeriodCredit:      dto.PeriodCredit,
		EndingAmount:      dto.EndingAmount,
	}
}

func ledgerSummaryToVO(dto query.LedgerSummary) LedgerSummaryResponse {
	return LedgerSummaryResponse(dto)
}

func auxiliaryLedgerSummaryToVO(dto query.AuxiliaryLedgerSummary) AuxiliaryLedgerSummaryResponse {
	return AuxiliaryLedgerSummaryResponse(dto)
}

func journalLineDTOToVO(dto query.JournalLine) JournalLineResponse {
	return JournalLineResponse{
		Id:                dto.Id,
		Account:           accountDTOToVO(dto.Account),
		AuxiliaryAccounts: converter.DTOsToVOs(dto.AuxiliaryAccounts, auxiliaryAccountDTOToVO),
		Text:              dto.Text,
		Amount:            dto.Amount,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
	}
}

func journalDTOToVO(dto query.Journal) JournalResponse {
	userOrNil := func(u *query.User) *UserResponse {
		if u != nil {
			return &UserResponse{
				Id:     u.Id,
				Traits: u.Traits,
			}
		}
		return nil
	}

	return JournalResponse{
		SobId:              dto.SobId,
		Id:                 dto.Id,
		Period:             periodDTOToVO(dto.Period),
		HeaderText:         dto.HeaderText,
		JournalType:        dto.JournalType,
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
		JournalLines:       converter.DTOsToVOs(dto.JournalLines, journalLineDTOToVO),
		CreatedAt:          dto.CreatedAt,
		UpdatedAt:          dto.UpdatedAt,
	}
}

func ledgerEntryDTOToVO(dto query.LedgerEntry) LedgerEntryResponse {
	return LedgerEntryResponse{
		JournalId:       dto.JournalId,
		JournalNumber:   dto.JournalNumber,
		TransactionDate: dto.TransactionDate,
		Text:            dto.Text,
		Amount:          dto.Amount,
		CreatedAt:       dto.CreatedAt,
		UpdatedAt:       dto.UpdatedAt,
	}
}
