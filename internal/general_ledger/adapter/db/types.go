package db

import (
	"fmt"
	"strings"
	"time"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/shopspring/decimal"
)

type accountPO struct {
	Id                uuid.UUID  `gorm:"type:uuid;primaryKey"`
	SobId             uuid.UUID  `gorm:"type:uuid;uniqueIndex:UQ_Accounts_SobId_AccountNumber"`
	SuperiorAccountId *uuid.UUID `gorm:"type:uuid"`
	Title             string
	AccountNumber     string           `gorm:"uniqueIndex:UQ_Accounts_SobId_AccountNumber"`
	NumberHierarchy   pgtype.Int4Array `gorm:"type:integer[]"`
	Level             int
	IsLeaf            bool
	Class             int
	Group             int
	BalanceDirection  string

	DimensionCategories []accountDimensionCategoryPO `gorm:"foreignKey:AccountId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// accountDimensionCategoryPO is the join table linking accounts to their allowed dimension categories.
type accountDimensionCategoryPO struct {
	AccountId           uuid.UUID `gorm:"type:uuid;primaryKey"`
	DimensionCategoryId uuid.UUID `gorm:"type:uuid;primaryKey"`
}

type periodPO struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId        uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_Periods_SobId_FiscalYear_PeriodNumber"`
	FiscalYear   int       `gorm:"uniqueIndex:UQ_Periods_SobId_FiscalYear_PeriodNumber"`
	PeriodNumber int       `gorm:"uniqueIndex:UQ_Periods_SobId_FiscalYear_PeriodNumber"`
	IsClosed     bool
	IsCurrent    bool

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type ledgerPO struct {
	Id            uuid.UUID       `gorm:"type:uuid;primaryKey"`
	SobId         uuid.UUID       `gorm:"type:uuid"`
	AccountId     uuid.UUID       `gorm:"type:uuid"`
	PeriodId      uuid.UUID       `gorm:"type:uuid"`
	OpeningAmount decimal.Decimal `gorm:"type:numeric"`
	PeriodAmount  decimal.Decimal `gorm:"type:numeric"`
	PeriodDebit   decimal.Decimal `gorm:"type:numeric"`
	PeriodCredit  decimal.Decimal `gorm:"type:numeric"`
	EndingAmount  decimal.Decimal `gorm:"type:numeric"`

	Account accountPO `gorm:"foreignKey:AccountId"`
	Period  periodPO  `gorm:"foreignKey:PeriodId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type journalPO struct {
	Id                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId              uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_Journals_SobId_PeriodId_DocumentNumber"`
	PeriodId           uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_Journals_SobId_PeriodId_DocumentNumber"`
	JournalType        string
	HeaderText         string
	DocumentNumber     string `gorm:"uniqueIndex:UQ_Journals_SobId_PeriodId_DocumentNumber"`
	AttachmentQuantity int
	Amount             decimal.Decimal `gorm:"type:numeric"`
	Creator            uuid.UUID       `gorm:"type:uuid"`
	Reviewer           uuid.UUID       `gorm:"type:uuid"`
	Auditor            uuid.UUID       `gorm:"type:uuid"`
	Poster             uuid.UUID       `gorm:"type:uuid"`
	IsReviewed         bool
	IsAudited          bool
	IsPosted           bool
	TransactionDate    time.Time `gorm:"type:date"`

	JournalLines []journalLinePO `gorm:"foreignKey:JournalId"`
	Period       periodPO        `gorm:"foreignKey:PeriodId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type journalLinePO struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	JournalId uuid.UUID `gorm:"type:uuid"`
	AccountId uuid.UUID `gorm:"type:uuid"`
	Text      string
	Amount    decimal.Decimal `gorm:"type:numeric"`

	Journal          journalPO                      `gorm:"foreignKey:JournalId"`
	Account          accountPO                      `gorm:"foreignKey:AccountId"`
	DimensionOptions []journalLineDimensionOptionPO `gorm:"foreignKey:JournalLineId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// journalLineDimensionOptionPO is the join table linking journal lines to their dimension options.
type journalLineDimensionOptionPO struct {
	JournalLineId     uuid.UUID `gorm:"type:uuid;primaryKey"`
	DimensionOptionId uuid.UUID `gorm:"type:uuid;primaryKey"`
}

// schemas

func (a accountPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return "accounts", nil
	}
	return "", fmt.Errorf("accountPO doesn't have association named %s", entity)
}

func (p periodPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return "periods", nil
	}
	return "", fmt.Errorf("periodPO doesn't have association named %s", entity)
}

func (l ledgerPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return "ledgers", nil
	}
	if strings.EqualFold(entity, "account") {
		return "Account", nil
	}
	return "", fmt.Errorf("ledgerPO doesn't have association named %s", entity)
}

func (j journalPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return "journals", nil
	}
	if strings.EqualFold(entity, "journalLines") {
		return "JournalLines", nil
	}
	if strings.EqualFold(entity, "period") {
		return "Period", nil
	}
	return "", fmt.Errorf("journalPO doesn't have association named %s", entity)
}

func (j journalLinePO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return "journal_lines", nil
	}
	if strings.EqualFold(entity, "account") {
		return "Account", nil
	}
	if strings.EqualFold(entity, "journal") {
		return "Journal", nil
	}
	return "", fmt.Errorf("journalLinePO doesn't have association named %s", entity)
}

// mappers

func accountBOToPO(bo *account.Account) accountPO {
	var int4array pgtype.Int4Array
	if err := int4array.Set(bo.NumberHierarchy()); err != nil {
		panic(fmt.Errorf("failde to convert []int to Int4Array: %w", err))
	}

	dimCategories := make([]accountDimensionCategoryPO, 0, len(bo.DimensionCategoryIds()))
	for _, catId := range bo.DimensionCategoryIds() {
		dimCategories = append(dimCategories, accountDimensionCategoryPO{
			AccountId:           bo.Id(),
			DimensionCategoryId: catId,
		})
	}

	return accountPO{
		Id:                  bo.Id(),
		SobId:               bo.SobId(),
		SuperiorAccountId:   converter.UUIDToPtr(bo.SuperiorAccountId()),
		Title:               bo.Title(),
		AccountNumber:       bo.AccountNumber(),
		NumberHierarchy:     int4array,
		Level:               bo.Level(),
		IsLeaf:              bo.IsLeaf(),
		Class:               int(bo.Class()),
		Group:               int(bo.Group()),
		BalanceDirection:    bo.BalanceDirection().String(),
		DimensionCategories: dimCategories,
	}
}

func accountPOToBO(po accountPO) (*account.Account, error) {
	var numberHierarchy []int
	if err := po.NumberHierarchy.AssignTo(&numberHierarchy); err != nil {
		return nil, fmt.Errorf("failed to assign Int4Array to []int: %w", err)
	}

	dimCategoryIds := make([]uuid.UUID, 0, len(po.DimensionCategories))
	for _, dc := range po.DimensionCategories {
		dimCategoryIds = append(dimCategoryIds, dc.DimensionCategoryId)
	}

	return account.NewByAllFields(
		po.Id,
		po.SobId,
		converter.UUIDFromPtr(po.SuperiorAccountId),
		nil,
		po.Title,
		po.AccountNumber,
		numberHierarchy,
		po.Level,
		po.IsLeaf,
		po.Class,
		po.Group,
		po.BalanceDirection,
		dimCategoryIds,
	)
}

func accountPOToBOWithSuperior(po accountPO, superior *account.Account) (*account.Account, error) {
	var numberHierarchy []int
	if err := po.NumberHierarchy.AssignTo(&numberHierarchy); err != nil {
		return nil, fmt.Errorf("failed to assign Int4Array to []int: %w", err)
	}

	dimCategoryIds := make([]uuid.UUID, 0, len(po.DimensionCategories))
	for _, dc := range po.DimensionCategories {
		dimCategoryIds = append(dimCategoryIds, dc.DimensionCategoryId)
	}

	return account.NewByAllFields(
		po.Id,
		po.SobId,
		converter.UUIDFromPtr(po.SuperiorAccountId),
		superior,
		po.Title,
		po.AccountNumber,
		numberHierarchy,
		po.Level,
		po.IsLeaf,
		po.Class,
		po.Group,
		po.BalanceDirection,
		dimCategoryIds,
	)
}

func accountPOToDTO(po accountPO) query.Account {
	var numberHierarchy []int
	if err := po.NumberHierarchy.AssignTo(&numberHierarchy); err != nil {
		panic(fmt.Errorf("failed to assign Int4Array to []int: %w", err))
	}

	dimCategoryIds := make([]uuid.UUID, 0, len(po.DimensionCategories))
	for _, dc := range po.DimensionCategories {
		dimCategoryIds = append(dimCategoryIds, dc.DimensionCategoryId)
	}

	return query.Account{
		SobId:                po.SobId,
		Id:                   po.Id,
		SuperiorAccountId:    po.SuperiorAccountId,
		Title:                po.Title,
		AccountNumber:        po.AccountNumber,
		NumberHierarchy:      numberHierarchy,
		Level:                po.Level,
		IsLeaf:               po.IsLeaf,
		Class:                po.Class,
		Group:                po.Group,
		BalanceDirection:     po.BalanceDirection,
		DimensionCategoryIds: dimCategoryIds,
		CreatedAt:            po.CreatedAt,
		UpdatedAt:            po.UpdatedAt,
	}
}

func periodBOToPO(bo period.Period) periodPO {
	return periodPO{
		SobId:        bo.SobId(),
		Id:           bo.Id(),
		FiscalYear:   bo.FiscalYear(),
		PeriodNumber: bo.PeriodNumber(),
		IsClosed:     bo.IsClosed(),
		IsCurrent:    bo.IsCurrent(),
	}
}

func periodPOToBO(po periodPO) (*period.Period, error) {
	return period.NewByAllFields(
		po.Id,
		po.SobId,
		po.FiscalYear,
		po.PeriodNumber,
		po.IsClosed,
		po.IsCurrent,
	)
}

func periodPOToDTO(po periodPO) query.Period {
	return query.Period(po)
}

func ledgerBOToPO(bo *ledger.Ledger) ledgerPO {
	return ledgerPO{
		Id:            bo.Id(),
		SobId:         bo.SobId(),
		AccountId:     bo.AccountId(),
		PeriodId:      bo.PeriodId(),
		OpeningAmount: bo.OpeningAmount(),
		PeriodAmount:  bo.PeriodAmount(),
		PeriodDebit:   bo.PeriodDebit(),
		PeriodCredit:  bo.PeriodCredit(),
		EndingAmount:  bo.EndingAmount(),
	}
}

func ledgerPOToBO(po ledgerPO) (*ledger.Ledger, error) {
	accountBO, err := accountPOToBO(po.Account)
	if err != nil {
		return nil, err
	}

	return ledger.New(
		po.Id,
		po.SobId,
		po.PeriodId,
		po.AccountId,
		accountBO,
		po.OpeningAmount,
		po.PeriodAmount,
		po.PeriodDebit,
		po.PeriodCredit,
		po.EndingAmount,
	)
}

func ledgerPOToDTO(po ledgerPO) query.Ledger {
	accountDTO := accountPOToDTO(po.Account)

	return query.Ledger{
		Id:            po.Id,
		SobId:         po.SobId,
		AccountId:     po.AccountId,
		PeriodId:      po.PeriodId,
		OpeningAmount: po.OpeningAmount,
		PeriodAmount:  po.PeriodAmount,
		PeriodDebit:   po.PeriodDebit,
		PeriodCredit:  po.PeriodCredit,
		EndingAmount:  po.EndingAmount,
		Account:       accountDTO,
		CreatedAt:     po.CreatedAt,
		UpdatedAt:     po.UpdatedAt,
	}
}

func journalBOToPO(bo journal.Journal) journalPO {
	var linePOs []journalLinePO
	for _, line := range bo.JournalLines() {
		linePOs = append(linePOs, journalLineBOToPO(*line, bo.Id()))
	}

	// Convert TransactionDate to time.Time for PostgreSQL DATE type
	transactionDate := time.Date(
		bo.TransactionDate().Year,
		time.Month(bo.TransactionDate().Month),
		bo.TransactionDate().Day,
		0, 0, 0, 0, time.UTC,
	)

	return journalPO{
		SobId:              bo.SobId(),
		Id:                 bo.Id(),
		PeriodId:           bo.PeriodId(),
		JournalType:        bo.JournalType().String(),
		HeaderText:         bo.HeaderText(),
		DocumentNumber:     bo.DocumentNumber(),
		AttachmentQuantity: bo.AttachmentQuantity(),
		Amount:             bo.Amount(),
		Creator:            bo.Creator(),
		Reviewer:           bo.Reviewer(),
		Auditor:            bo.Auditor(),
		Poster:             bo.Poster(),
		IsReviewed:         bo.IsReviewed(),
		IsAudited:          bo.IsAudited(),
		IsPosted:           bo.IsPosted(),
		TransactionDate:    transactionDate,
		JournalLines:       linePOs,
	}
}

func journalPOToBO(po journalPO) (*journal.Journal, error) {
	lineBOs, err := converter.POsToBOs(po.JournalLines, journalLinePOToBO)
	if err != nil {
		return nil, err
	}

	periodBO, err := periodPOToBO(po.Period)
	if err != nil {
		return nil, err
	}

	// Convert time.Time DATE to TransactionDate
	transactionDate := transaction_date.TransactionDate{
		Year:  po.TransactionDate.Year(),
		Month: int(po.TransactionDate.Month()),
		Day:   po.TransactionDate.Day(),
	}

	return journal.New(
		po.Id,
		po.SobId,
		periodBO,
		po.JournalType,
		po.HeaderText,
		po.DocumentNumber,
		po.AttachmentQuantity,
		po.Creator,
		po.Reviewer,
		po.Auditor,
		po.Poster,
		po.IsReviewed,
		po.IsAudited,
		po.IsPosted,
		transactionDate,
		lineBOs,
	)
}

func journalPOToDTO(po journalPO) query.Journal {
	periodDTO := periodPOToDTO(po.Period)

	lineDTOs := converter.POsToDTOs(po.JournalLines, journalLinePOToDTO)

	userOrNil := func(id uuid.UUID) *query.User {
		if id != uuid.Nil {
			return &query.User{Id: id}
		}
		return nil
	}

	// Convert time.Time DATE to TransactionDate
	transactionDate := transaction_date.TransactionDate{
		Year:  po.TransactionDate.Year(),
		Month: int(po.TransactionDate.Month()),
		Day:   po.TransactionDate.Day(),
	}

	return query.Journal{
		SobId:              po.SobId,
		Id:                 po.Id,
		Period:             periodDTO,
		JournalType:        po.JournalType,
		HeaderText:         po.HeaderText,
		DocumentNumber:     po.DocumentNumber,
		AttachmentQuantity: po.AttachmentQuantity,
		Amount:             po.Amount,
		Creator:            userOrNil(po.Creator),
		Reviewer:           userOrNil(po.Reviewer),
		Auditor:            userOrNil(po.Auditor),
		Poster:             userOrNil(po.Poster),
		IsReviewed:         po.IsReviewed,
		IsAudited:          po.IsAudited,
		IsPosted:           po.IsPosted,
		TransactionDate:    transactionDate,
		JournalLines:       lineDTOs,
		CreatedAt:          po.CreatedAt,
		UpdatedAt:          po.UpdatedAt,
	}
}

func journalLineBOToPO(bo journal.JournalLine, journalId uuid.UUID) journalLinePO {
	dimOptions := make([]journalLineDimensionOptionPO, 0, len(bo.DimensionOptionIds()))
	for _, optId := range bo.DimensionOptionIds() {
		dimOptions = append(dimOptions, journalLineDimensionOptionPO{
			JournalLineId:     bo.Id(),
			DimensionOptionId: optId,
		})
	}

	return journalLinePO{
		JournalId:        journalId,
		Id:               bo.Id(),
		AccountId:        bo.AccountId(),
		Text:             bo.Text(),
		Amount:           bo.Amount(),
		DimensionOptions: dimOptions,
	}
}

func journalLinePOToBO(po journalLinePO) (*journal.JournalLine, error) {
	accountBO, err := accountPOToBO(po.Account)
	if err != nil {
		return nil, err
	}

	dimOptionIds := make([]uuid.UUID, 0, len(po.DimensionOptions))
	for _, d := range po.DimensionOptions {
		dimOptionIds = append(dimOptionIds, d.DimensionOptionId)
	}

	return journal.NewJournalLine(
		po.Id,
		accountBO,
		po.Text,
		po.Amount,
		dimOptionIds,
	)
}

func journalLinePOToDTO(po journalLinePO) query.JournalLine {
	dimOptionIds := make([]uuid.UUID, 0, len(po.DimensionOptions))
	for _, d := range po.DimensionOptions {
		dimOptionIds = append(dimOptionIds, d.DimensionOptionId)
	}

	return query.JournalLine{
		Id:                 po.Id,
		Account:            accountPOToDTO(po.Account),
		Text:               po.Text,
		Amount:             po.Amount,
		DimensionOptionIds: dimOptionIds,
		CreatedAt:          po.CreatedAt,
		UpdatedAt:          po.UpdatedAt,
	}
}
