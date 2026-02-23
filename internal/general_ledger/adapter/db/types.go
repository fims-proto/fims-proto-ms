package db

import (
	"fmt"
	"strings"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_category"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

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

	AuxiliaryCategories []auxiliaryCategoryPO `gorm:"many2many:account_auxiliary_category_links;joinForeignKey:AccountId;joinReferences:AuxiliaryCategoryId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type auxiliaryCategoryPO struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId      uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_AuxiliaryCategories_SobId_Key;uniqueIndex:UQ_AuxiliaryCategories_SobId_Title"`
	Key        string    `gorm:"uniqueIndex:UQ_AuxiliaryCategories_SobId_Key"`
	Title      string    `gorm:"uniqueIndex:UQ_AuxiliaryCategories_SobId_Title"`
	IsStandard bool

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type auxiliaryAccountPO struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey"`
	CategoryId  uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_AuxiliaryAccounts_CategoryId_Key;uniqueIndex:UQ_AuxiliaryAccounts_CategoryId_Title"`
	Key         string    `gorm:"uniqueIndex:UQ_AuxiliaryAccounts_CategoryId_Key"`
	Title       string    `gorm:"uniqueIndex:UQ_AuxiliaryAccounts_CategoryId_Title"`
	Description string

	Category auxiliaryCategoryPO `gorm:"foreignKey:CategoryId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
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
	Id                   uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId                uuid.UUID `gorm:"type:uuid"`
	AccountId            uuid.UUID `gorm:"type:uuid"`
	PeriodId             uuid.UUID `gorm:"type:uuid"`
	OpeningDebitBalance  decimal.Decimal
	OpeningCreditBalance decimal.Decimal
	PeriodDebit          decimal.Decimal
	PeriodCredit         decimal.Decimal
	EndingDebitBalance   decimal.Decimal
	EndingCreditBalance  decimal.Decimal

	Account accountPO `gorm:"foreignKey:AccountId"`
	Period  periodPO  `gorm:"foreignKey:PeriodId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type auxiliaryLedgerPO struct {
	Id                   uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId                uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_auxiliary_ledgers_natural_key"`
	PeriodId             uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_auxiliary_ledgers_natural_key"`
	AccountId            uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_auxiliary_ledgers_natural_key"`
	AuxiliaryCategoryId  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_auxiliary_ledgers_natural_key"`
	AuxiliaryAccountId   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:uq_auxiliary_ledgers_natural_key"`
	OpeningDebitBalance  decimal.Decimal
	OpeningCreditBalance decimal.Decimal
	PeriodDebit          decimal.Decimal
	PeriodCredit         decimal.Decimal
	EndingDebitBalance   decimal.Decimal
	EndingCreditBalance  decimal.Decimal

	Account           accountPO           `gorm:"foreignKey:AccountId"`
	AuxiliaryCategory auxiliaryCategoryPO `gorm:"foreignKey:AuxiliaryCategoryId"`
	AuxiliaryAccount  auxiliaryAccountPO  `gorm:"foreignKey:AuxiliaryAccountId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type voucherPO struct {
	Id                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId              uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_Vouchers_SobId_PeriodId_DocumentNumber"`
	PeriodId           uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_Vouchers_SobId_PeriodId_DocumentNumber"`
	VoucherType        string
	HeaderText         string
	DocumentNumber     string `gorm:"uniqueIndex:UQ_Vouchers_SobId_PeriodId_DocumentNumber"`
	AttachmentQuantity int
	Debit              decimal.Decimal
	Credit             decimal.Decimal
	Creator            uuid.UUID `gorm:"type:uuid"`
	Reviewer           uuid.UUID `gorm:"type:uuid"`
	Auditor            uuid.UUID `gorm:"type:uuid"`
	Poster             uuid.UUID `gorm:"type:uuid"`
	IsReviewed         bool
	IsAudited          bool
	IsPosted           bool
	TransactionDate    time.Time `gorm:"type:date"`

	LineItems []lineItemPO `gorm:"foreignKey:VoucherId"`
	Period    periodPO     `gorm:"foreignKey:PeriodId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type lineItemPO struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	VoucherId uuid.UUID `gorm:"type:uuid"`
	AccountId uuid.UUID `gorm:"type:uuid"`
	Text      string
	Debit     decimal.Decimal
	Credit    decimal.Decimal

	Account           accountPO            `gorm:"foreignKey:AccountId"`
	AuxiliaryAccounts []auxiliaryAccountPO `gorm:"many2many:line_item_auxiliary_account_links;joinForeignKey:LineItemId;joinReferences:AuxiliaryAccountId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// table names

func (a accountPO) TableName() string {
	return "a_accounts"
}

func (a auxiliaryCategoryPO) TableName() string {
	return "a_auxiliary_categories"
}

func (a auxiliaryAccountPO) TableName() string {
	return "a_auxiliary_accounts"
}

func (p periodPO) TableName() string {
	return "a_periods"
}

func (l ledgerPO) TableName() string {
	return "a_ledgers"
}

func (a auxiliaryLedgerPO) TableName() string {
	return "a_auxiliary_ledgers"
}

func (v voucherPO) TableName() string {
	return "a_vouchers"
}

func (l lineItemPO) TableName() string {
	return "a_line_items"
}

// schemas

func (a accountPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return a.TableName(), nil
	}
	if strings.EqualFold(entity, "auxiliaryCategories") {
		return "AuxiliaryCategories", nil
	}
	return "", fmt.Errorf("accountPO doesn't have association named %s", entity)
}

func (a auxiliaryCategoryPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return a.TableName(), nil
	}
	return "", fmt.Errorf("auxiliaryCategoryPO doesn't have association named %s", entity)
}

func (a auxiliaryAccountPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return a.TableName(), nil
	}
	if strings.EqualFold(entity, "category") {
		return "Category", nil
	}
	return "", fmt.Errorf("auxiliaryAccountPO doesn't have association named %s", entity)
}

func (p periodPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return p.TableName(), nil
	}
	return "", fmt.Errorf("periodPO doesn't have association named %s", entity)
}

func (l ledgerPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return l.TableName(), nil
	}
	if strings.EqualFold(entity, "account") {
		return "Account", nil
	}
	return "", fmt.Errorf("ledgerPO doesn't have association named %s", entity)
}

func (a auxiliaryLedgerPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return a.TableName(), nil
	}
	if strings.EqualFold(entity, "auxiliaryAccount") {
		return "AuxiliaryAccount", nil
	}
	return "", fmt.Errorf("auxiliaryLedgerPO doesn't have association named %s", entity)
}

func (v voucherPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return v.TableName(), nil
	}
	if strings.EqualFold(entity, "lineItems") {
		return "LineItems", nil
	}
	if strings.EqualFold(entity, "period") {
		return "Period", nil
	}
	return "", fmt.Errorf("voucherPO doesn't have association named %s", entity)
}

func (l lineItemPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return l.TableName(), nil
	}
	if strings.EqualFold(entity, "account") {
		return "Account", nil
	}
	if strings.EqualFold(entity, "auxiliaryAccount") {
		return "AuxiliaryAccount", nil
	}
	return "", fmt.Errorf("lineItemPO doesn't have association named %s", entity)
}

// mappers

func accountBOToPO(bo *account.Account) accountPO {
	var int4array pgtype.Int4Array
	if err := int4array.Set(bo.NumberHierarchy()); err != nil {
		panic(fmt.Errorf("failde to convert []int to Int4Array: %w", err))
	}

	var categoryPOs []auxiliaryCategoryPO
	for _, category := range bo.AuxiliaryCategories() {
		categoryPOs = append(categoryPOs, auxiliaryCategoryBOToPO(category))
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
		AuxiliaryCategories: categoryPOs,
	}
}

func accountPOToBO(po accountPO) (*account.Account, error) {
	var numberHierarchy []int
	if err := po.NumberHierarchy.AssignTo(&numberHierarchy); err != nil {
		return nil, fmt.Errorf("failed to assign Int4Array to []int: %w", err)
	}

	categoryBOs, err := converter.POsToBOs(po.AuxiliaryCategories, auxiliaryCategoryPOToBO)
	if err != nil {
		return nil, err
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
		categoryBOs,
	)
}

func accountPOToBOWithSuperior(po accountPO, superior *account.Account) (*account.Account, error) {
	var numberHierarchy []int
	if err := po.NumberHierarchy.AssignTo(&numberHierarchy); err != nil {
		return nil, fmt.Errorf("failed to assign Int4Array to []int: %w", err)
	}

	categoryBOs, err := converter.POsToBOs(po.AuxiliaryCategories, auxiliaryCategoryPOToBO)
	if err != nil {
		return nil, err
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
		categoryBOs,
	)
}

func accountPOToDTO(po accountPO) query.Account {
	var numberHierarchy []int
	if err := po.NumberHierarchy.AssignTo(&numberHierarchy); err != nil {
		panic(fmt.Errorf("failed to assign Int4Array to []int: %w", err))
	}

	categoryDTOs := converter.POsToDTOs(po.AuxiliaryCategories, auxiliaryCategoryPOToDTO)

	return query.Account{
		SobId:               po.SobId,
		Id:                  po.Id,
		SuperiorAccountId:   po.SuperiorAccountId,
		Title:               po.Title,
		AccountNumber:       po.AccountNumber,
		NumberHierarchy:     numberHierarchy,
		Level:               po.Level,
		IsLeaf:              po.IsLeaf,
		Class:               po.Class,
		Group:               po.Group,
		BalanceDirection:    po.BalanceDirection,
		AuxiliaryCategories: categoryDTOs,
		CreatedAt:           po.CreatedAt,
		UpdatedAt:           po.UpdatedAt,
	}
}

func auxiliaryCategoryPOToBO(po auxiliaryCategoryPO) (*auxiliary_category.AuxiliaryCategory, error) {
	return auxiliary_category.New(
		po.Id,
		po.SobId,
		po.Key,
		po.Title,
		po.IsStandard,
	)
}

func auxiliaryCategoryBOToPO(bo *auxiliary_category.AuxiliaryCategory) auxiliaryCategoryPO {
	return auxiliaryCategoryPO{
		Id:         bo.Id(),
		SobId:      bo.SobId(),
		Key:        bo.Key(),
		Title:      bo.Title(),
		IsStandard: bo.IsStandard(),
	}
}

func auxiliaryCategoryPOToDTO(po auxiliaryCategoryPO) query.AuxiliaryCategory {
	return query.AuxiliaryCategory{
		Id:         po.Id,
		SobId:      po.SobId,
		Key:        po.Key,
		Title:      po.Title,
		IsStandard: po.IsStandard,
		CreatedAt:  po.CreatedAt,
		UpdatedAt:  po.UpdatedAt,
	}
}

func auxiliaryAccountPOToBO(po auxiliaryAccountPO) (*auxiliary_account.AuxiliaryAccount, error) {
	categoryBO, err := auxiliaryCategoryPOToBO(po.Category)
	if err != nil {
		return nil, err
	}

	return auxiliary_account.New(
		po.Id,
		categoryBO,
		po.Key,
		po.Title,
		po.Description,
	)
}

func auxiliaryAccountBOToPO(bo *auxiliary_account.AuxiliaryAccount) auxiliaryAccountPO {
	return auxiliaryAccountPO{
		Id:          bo.Id(),
		CategoryId:  bo.Category().Id(),
		Key:         bo.Key(),
		Title:       bo.Title(),
		Description: bo.Description(),
		Category:    auxiliaryCategoryBOToPO(bo.Category()),
	}
}

func auxiliaryAccountPOToDTO(po auxiliaryAccountPO) query.AuxiliaryAccount {
	return query.AuxiliaryAccount{
		Id:          po.Id,
		Category:    auxiliaryCategoryPOToDTO(po.Category),
		Key:         po.Key,
		Title:       po.Title,
		Description: po.Description,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
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
		Id:                   bo.Id(),
		SobId:                bo.SobId(),
		AccountId:            bo.AccountId(),
		PeriodId:             bo.PeriodId(),
		OpeningDebitBalance:  bo.OpeningDebitBalance(),
		OpeningCreditBalance: bo.OpeningCreditBalance(),
		PeriodDebit:          bo.PeriodDebit(),
		PeriodCredit:         bo.PeriodCredit(),
		EndingDebitBalance:   bo.EndingDebitBalance(),
		EndingCreditBalance:  bo.EndingCreditBalance(),
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
		po.OpeningDebitBalance,
		po.OpeningCreditBalance,
		po.PeriodDebit,
		po.PeriodCredit,
		po.EndingDebitBalance,
		po.EndingCreditBalance,
	)
}

func ledgerPOToDTO(po ledgerPO) query.Ledger {
	accountDTO := accountPOToDTO(po.Account)

	return query.Ledger{
		Id:                   po.Id,
		SobId:                po.SobId,
		AccountId:            po.AccountId,
		PeriodId:             po.PeriodId,
		OpeningDebitBalance:  po.OpeningDebitBalance,
		OpeningCreditBalance: po.OpeningCreditBalance,
		PeriodDebit:          po.PeriodDebit,
		PeriodCredit:         po.PeriodCredit,
		EndingDebitBalance:   po.EndingDebitBalance,
		EndingCreditBalance:  po.EndingCreditBalance,
		Account:              accountDTO,
		CreatedAt:            po.CreatedAt,
		UpdatedAt:            po.UpdatedAt,
	}
}

func auxiliaryLedgerBOToPO(bo *auxiliary_ledger.AuxiliaryLedger) auxiliaryLedgerPO {
	return auxiliaryLedgerPO{
		Id:                   bo.Id(),
		SobId:                bo.SobId(),
		PeriodId:             bo.PeriodId(),
		AccountId:            bo.AccountId(),
		AuxiliaryCategoryId:  bo.AuxiliaryCategoryId(),
		AuxiliaryAccountId:   bo.AuxiliaryAccountId(),
		OpeningDebitBalance:  bo.OpeningDebitBalance(),
		OpeningCreditBalance: bo.OpeningCreditBalance(),
		PeriodDebit:          bo.PeriodDebit(),
		PeriodCredit:         bo.PeriodCredit(),
		EndingDebitBalance:   bo.EndingDebitBalance(),
		EndingCreditBalance:  bo.EndingCreditBalance(),
	}
}

func auxiliaryLedgerPOToBO(po auxiliaryLedgerPO) (*auxiliary_ledger.AuxiliaryLedger, error) {
	return auxiliary_ledger.New(
		po.Id,
		po.SobId,
		po.PeriodId,
		po.AccountId,
		po.AuxiliaryCategoryId,
		po.AuxiliaryAccountId,
		po.OpeningDebitBalance,
		po.OpeningCreditBalance,
		po.PeriodDebit,
		po.PeriodCredit,
		po.EndingDebitBalance,
		po.EndingCreditBalance,
	)
}

func auxiliaryLedgerPOToDTO(po auxiliaryLedgerPO) query.AuxiliaryLedger {
	return query.AuxiliaryLedger{
		Id:                   po.Id,
		SobId:                po.SobId,
		PeriodId:             po.PeriodId,
		Account:              accountPOToDTO(po.Account),
		AuxiliaryCategory:    auxiliaryCategoryPOToDTO(po.AuxiliaryCategory),
		AuxiliaryAccount:     auxiliaryAccountPOToDTO(po.AuxiliaryAccount),
		OpeningDebitBalance:  po.OpeningDebitBalance,
		OpeningCreditBalance: po.OpeningCreditBalance,
		PeriodDebit:          po.PeriodDebit,
		PeriodCredit:         po.PeriodCredit,
		EndingDebitBalance:   po.EndingDebitBalance,
		EndingCreditBalance:  po.EndingCreditBalance,
		CreatedAt:            po.CreatedAt,
		UpdatedAt:            po.UpdatedAt,
	}
}

func voucherBOToPO(bo voucher.Voucher) voucherPO {
	var itemPOs []lineItemPO
	for _, item := range bo.LineItems() {
		itemPOs = append(itemPOs, lineItemBOToPO(*item, bo.Id()))
	}

	// Convert TransactionDate to time.Time for PostgreSQL DATE type
	transactionDate := time.Date(
		bo.TransactionDate().Year,
		time.Month(bo.TransactionDate().Month),
		bo.TransactionDate().Day,
		0, 0, 0, 0, time.UTC,
	)

	return voucherPO{
		SobId:              bo.SobId(),
		Id:                 bo.Id(),
		PeriodId:           bo.PeriodId(),
		VoucherType:        bo.VoucherType().String(),
		HeaderText:         bo.HeaderText(),
		DocumentNumber:     bo.DocumentNumber(),
		AttachmentQuantity: bo.AttachmentQuantity(),
		Debit:              bo.Debit(),
		Credit:             bo.Credit(),
		Creator:            bo.Creator(),
		Reviewer:           bo.Reviewer(),
		Auditor:            bo.Auditor(),
		Poster:             bo.Poster(),
		IsReviewed:         bo.IsReviewed(),
		IsAudited:          bo.IsAudited(),
		IsPosted:           bo.IsPosted(),
		TransactionDate:    transactionDate,
		LineItems:          itemPOs,
	}
}

func voucherPOToBO(po voucherPO) (*voucher.Voucher, error) {
	itemBOs, err := converter.POsToBOs(po.LineItems, lineItemPOToBO)
	if err != nil {
		return nil, err
	}

	periodBO, err := periodPOToBO(po.Period)
	if err != nil {
		return nil, err
	}

	// Convert time.Time DATE to TransactionDate
	transactionDate := voucher.TransactionDate{
		Year:  po.TransactionDate.Year(),
		Month: int(po.TransactionDate.Month()),
		Day:   po.TransactionDate.Day(),
	}

	return voucher.New(
		po.Id,
		po.SobId,
		periodBO,
		po.VoucherType,
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
		itemBOs,
	)
}

func voucherPOToDTO(po voucherPO) query.Voucher {
	periodDTO := periodPOToDTO(po.Period)

	itemDTOs := converter.POsToDTOs(po.LineItems, lineItemPOToDTO)

	userOrNil := func(id uuid.UUID) *query.User {
		if id != uuid.Nil {
			return &query.User{Id: id}
		}
		return nil
	}

	// Convert time.Time DATE to TransactionDate
	transactionDate := voucher.TransactionDate{
		Year:  po.TransactionDate.Year(),
		Month: int(po.TransactionDate.Month()),
		Day:   po.TransactionDate.Day(),
	}

	return query.Voucher{
		SobId:              po.SobId,
		Id:                 po.Id,
		Period:             periodDTO,
		VoucherType:        po.VoucherType,
		HeaderText:         po.HeaderText,
		DocumentNumber:     po.DocumentNumber,
		AttachmentQuantity: po.AttachmentQuantity,
		Debit:              po.Debit,
		Credit:             po.Credit,
		Creator:            userOrNil(po.Creator),
		Reviewer:           userOrNil(po.Reviewer),
		Auditor:            userOrNil(po.Auditor),
		Poster:             userOrNil(po.Poster),
		IsReviewed:         po.IsReviewed,
		IsAudited:          po.IsAudited,
		IsPosted:           po.IsPosted,
		TransactionDate:    transactionDate,
		LineItems:          itemDTOs,
		CreatedAt:          po.CreatedAt,
		UpdatedAt:          po.UpdatedAt,
	}
}

func lineItemBOToPO(bo voucher.LineItem, voucherId uuid.UUID) lineItemPO {
	var auxiliaryAccounts []auxiliaryAccountPO
	for _, auxiliaryAccount := range bo.AuxiliaryAccounts() {
		auxiliaryAccounts = append(auxiliaryAccounts, auxiliaryAccountBOToPO(auxiliaryAccount))
	}

	return lineItemPO{
		VoucherId:         voucherId,
		Id:                bo.Id(),
		AccountId:         bo.AccountId(),
		AuxiliaryAccounts: auxiliaryAccounts,
		Text:              bo.Text(),
		Debit:             bo.Debit(),
		Credit:            bo.Credit(),
	}
}

func lineItemPOToBO(po lineItemPO) (*voucher.LineItem, error) {
	accountBO, err := accountPOToBO(po.Account)
	if err != nil {
		return nil, err
	}

	auxiliaryAccountBOs, err := converter.POsToBOs(po.AuxiliaryAccounts, auxiliaryAccountPOToBO)
	if err != nil {
		return nil, err
	}

	return voucher.NewLineItem(
		po.Id,
		accountBO,
		auxiliaryAccountBOs,
		po.Text,
		po.Debit,
		po.Credit,
	)
}

func lineItemPOToDTO(po lineItemPO) query.LineItem {
	accountDTO := accountPOToDTO(po.Account)

	var auxiliaryAccounts []query.AuxiliaryAccount
	for _, auxiliaryAccount := range po.AuxiliaryAccounts {
		auxiliaryAccounts = append(auxiliaryAccounts, auxiliaryAccountPOToDTO(auxiliaryAccount))
	}

	return query.LineItem{
		Id:                po.Id,
		Account:           accountDTO,
		AuxiliaryAccounts: auxiliaryAccounts,
		Text:              po.Text,
		Debit:             po.Debit,
		Credit:            po.Credit,
		CreatedAt:         po.CreatedAt,
		UpdatedAt:         po.UpdatedAt,
	}
}
