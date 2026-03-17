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

// DimensionCategoryResponse is embedded in AccountDetailResponse and DimensionOptionResponse.
type DimensionCategoryResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// DimensionOptionResponse is embedded in JournalLineResponse.
// Category is nested to give the full context without an extra round-trip.
type DimensionOptionResponse struct {
	Id       uuid.UUID                 `json:"id"`
	Name     string                    `json:"name"`
	Category DimensionCategoryResponse `json:"category"`
}

// AccountSlimResponse is used by list endpoints (GET /accounts).
// It only contains fields from the account table itself — no cross-table dimension data.
type AccountSlimResponse struct {
	Id                uuid.UUID  `json:"id,omitempty"`
	SobId             uuid.UUID  `json:"sobId,omitempty"`
	SuperiorAccountId *uuid.UUID `json:"superiorAccountId,omitempty"`
	Title             string     `json:"title,omitempty"`
	AccountNumber     string     `json:"accountNumber,omitempty"`
	NumberHierarchy   []int      `json:"numberHierarchy,omitempty"`
	Level             int        `json:"level"`
	IsLeaf            bool       `json:"isLeaf"`
	Class             string     `json:"class"`
	Group             string     `json:"group"`
	BalanceDirection  string     `json:"balanceDirection,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

// AccountDetailResponse is used by detail and create endpoints (GET /account/{id}, POST /accounts).
// It includes full dimension category objects.
type AccountDetailResponse struct {
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
	DimensionCategories []DimensionCategoryResponse `json:"dimensionCategories"`
	CreatedAt           time.Time                   `json:"createdAt"`
	UpdatedAt           time.Time                   `json:"updatedAt"`
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

type PeriodAndLedgersResponse struct {
	Period  PeriodResponse   `json:"period"`
	Ledgers []LedgerResponse `json:"ledgers"`
}

type JournalLineResponse struct {
	Id               uuid.UUID                 `json:"id,omitempty"`
	Account          AccountDetailResponse     `json:"account"`
	Text             string                    `json:"text,omitempty"`
	Amount           decimal.Decimal           `json:"amount"`
	DimensionOptions []DimensionOptionResponse `json:"dimensionOptions,omitempty"`
	CreatedAt        time.Time                 `json:"createdAt"`
	UpdatedAt        time.Time                 `json:"updatedAt"`
}

type JournalResponse struct {
	Id                 uuid.UUID                        `json:"id,omitempty"`
	SobId              uuid.UUID                        `json:"sobId,omitempty"`
	Period             PeriodResponse                   `json:"period"`
	HeaderText         string                           `json:"headerText,omitempty"`
	DocumentNumber     string                           `json:"documentNumber,omitempty"`
	JournalType        string                           `json:"journalType" enums:"GENERAL,ADJUSTING,REVERSING,CLOSING"`
	ReferenceJournalId *uuid.UUID                       `json:"referenceJournalId,omitempty"`
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

type LedgerDimensionSummaryItemResponse struct {
	DimensionOption DimensionOptionResponse `json:"dimensionOption"`
	TotalAmount     decimal.Decimal         `json:"totalAmount"`
}

// mapper

func accountDTOToSlimVO(dto query.Account) AccountSlimResponse {
	return AccountSlimResponse{
		Id:                dto.Id,
		SobId:             dto.SobId,
		SuperiorAccountId: dto.SuperiorAccountId,
		Title:             dto.Title,
		AccountNumber:     dto.AccountNumber,
		NumberHierarchy:   dto.NumberHierarchy,
		Level:             dto.Level,
		IsLeaf:            dto.IsLeaf,
		Class:             strconv.Itoa(dto.Class),
		Group:             strconv.Itoa(dto.Group),
		BalanceDirection:  dto.BalanceDirection,
		CreatedAt:         dto.CreatedAt,
		UpdatedAt:         dto.UpdatedAt,
	}
}

func accountDTOToDetailVO(dto query.Account) AccountDetailResponse {
	categories := make([]DimensionCategoryResponse, 0, len(dto.DimensionCategories))
	for _, cat := range dto.DimensionCategories {
		categories = append(categories, DimensionCategoryResponse{Id: cat.Id, Name: cat.Name})
	}

	return AccountDetailResponse{
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
		DimensionCategories: categories,
		CreatedAt:           dto.CreatedAt,
		UpdatedAt:           dto.UpdatedAt,
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

func journalLineDTOToVO(dto query.JournalLine) JournalLineResponse {
	options := make([]DimensionOptionResponse, 0, len(dto.DimensionOptions))
	for _, opt := range dto.DimensionOptions {
		options = append(options, DimensionOptionResponse{
			Id:   opt.Id,
			Name: opt.Name,
			Category: DimensionCategoryResponse{
				Id:   opt.Category.Id,
				Name: opt.Category.Name,
			},
		})
	}

	return JournalLineResponse{
		Id:               dto.Id,
		Account:          accountDTOToDetailVO(dto.Account),
		Text:             dto.Text,
		Amount:           dto.Amount,
		DimensionOptions: options,
		CreatedAt:        dto.CreatedAt,
		UpdatedAt:        dto.UpdatedAt,
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
		DocumentNumber:     dto.DocumentNumber,
		JournalType:        dto.JournalType,
		ReferenceJournalId: dto.ReferenceJournalId,
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

func ledgerDimensionSummaryItemToVO(dto query.LedgerDimensionSummaryItem) LedgerDimensionSummaryItemResponse {
	return LedgerDimensionSummaryItemResponse{
		DimensionOption: DimensionOptionResponse{
			Id:   dto.DimensionOptionId,
			Name: dto.DimensionOptionName,
		},
		TotalAmount: dto.TotalAmount,
	}
}
