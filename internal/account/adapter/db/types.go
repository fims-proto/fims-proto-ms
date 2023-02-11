package db

import (
	"strings"
	"time"

	"github/fims-proto/fims-proto-ms/internal/account/domain/ledger"

	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"github/fims-proto/fims-proto-ms/internal/account/domain/period"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

type accountPO struct {
	Id                uuid.UUID `gorm:"type:uuid"`
	SobId             uuid.UUID `gorm:"type:uuid;uniqueIndex:accounts_sobid_number_key"`
	SuperiorAccountId uuid.UUID `gorm:"type:uuid"`
	Title             string
	AccountNumber     string           `gorm:"uniqueIndex:accounts_sobid_number_key"`
	NumberHierarchy   pgtype.Int4Array `gorm:"type:integer[]"`
	Level             int
	AccountType       string
	BalanceDirection  string
	CreatedAt         time.Time `gorm:"<-:create"`
	UpdatedAt         time.Time
}

type periodPO struct {
	Id           uuid.UUID `gorm:"type:uuid"`
	SobId        uuid.UUID `gorm:"type:uuid;uniqueIndex:periods_sobid_year_number_key"`
	FiscalYear   int       `gorm:"uniqueIndex:periods_sobid_year_number_key"`
	PeriodNumber int       `gorm:"uniqueIndex:periods_sobid_year_number_key"`
	OpeningTime  time.Time
	EndingTime   time.Time
	IsClosed     bool
	CreatedAt    time.Time `gorm:"<-:create"`
	UpdatedAt    time.Time
}

type ledgerPO struct {
	Id             uuid.UUID `gorm:"type:uuid"`
	SobId          uuid.UUID `gorm:"type:uuid"`
	AccountId      uuid.UUID `gorm:"type:uuid;primaryKey"`
	PeriodId       uuid.UUID `gorm:"type:uuid;primaryKey"`
	OpeningBalance decimal.Decimal
	EndingBalance  decimal.Decimal
	PeriodDebit    decimal.Decimal
	PeriodCredit   decimal.Decimal
	Account        accountPO `gorm:"foreignKey:AccountId"`
	CreatedAt      time.Time `gorm:"<-:create"`
	UpdatedAt      time.Time
}

// table names

func (a accountPO) TableName() string {
	return "a_accounts"
}

func (p periodPO) TableName() string {
	return "a_periods"
}

func (l ledgerPO) TableName() string {
	return "a_ledgers"
}

// schemas

func (a accountPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return a.TableName(), nil
	}
	return "", errors.Errorf("accountPO doesn't have association named %s", entity)
}

func (p periodPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return p.TableName(), nil
	}
	return "", errors.Errorf("periodPO doesn't have association named %s", entity)
}

func (l ledgerPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return l.TableName(), nil
	}
	if strings.EqualFold(entity, "account") {
		return "Account", nil
	}
	return "", errors.Errorf("ledgerPO doesn't have association named %s", entity)
}

// mappers

func accountBOToPO(bo account.Account) (accountPO, error) {
	var int4array pgtype.Int4Array
	if err := int4array.Set(bo.NumberHierarchy()); err != nil {
		return accountPO{}, errors.Wrap(err, "convert []int to Int4Array failed")
	}
	return accountPO{
		Id:                bo.Id(),
		SobId:             bo.SobId(),
		SuperiorAccountId: bo.SuperiorAccountId(),
		Title:             bo.Title(),
		AccountNumber:     bo.AccountNumber(),
		NumberHierarchy:   int4array,
		Level:             bo.Level(),
		AccountType:       bo.AccountType().String(),
		BalanceDirection:  bo.BalanceDirection().String(),
	}, nil
}

func accountPOToBO(po accountPO) (*account.Account, error) {
	var numberHierarchy []int
	if err := po.NumberHierarchy.AssignTo(&numberHierarchy); err != nil {
		return nil, errors.Wrap(err, "assign Int4Array to []int failed")
	}

	return account.New(po.Id, po.SobId, po.SuperiorAccountId, po.Title, po.AccountNumber, numberHierarchy, po.Level, po.AccountType, po.BalanceDirection)
}

func accountPOToDTO(po accountPO) (query.Account, error) {
	var numberHierarchy []int
	if err := po.NumberHierarchy.AssignTo(&numberHierarchy); err != nil {
		return query.Account{}, errors.Wrap(err, "assign Int4Array to []int failed")
	}

	return query.Account{
		SobId:             po.SobId,
		Id:                po.Id,
		SuperiorAccountId: po.SuperiorAccountId,
		Title:             po.Title,
		AccountNumber:     po.AccountNumber,
		NumberHierarchy:   numberHierarchy,
		Level:             po.Level,
		AccountType:       po.AccountType,
		BalanceDirection:  po.BalanceDirection,
	}, nil
}

func periodBOToPO(bo period.Period) periodPO {
	return periodPO{
		SobId:        bo.SobId(),
		Id:           bo.Id(),
		FiscalYear:   bo.FiscalYear(),
		PeriodNumber: bo.PeriodNumber(),
		OpeningTime:  bo.OpeningTime(),
		EndingTime:   bo.EndingTime(),
		IsClosed:     bo.IsClosed(),
	}
}

func periodPOToDTO(po periodPO) (query.Period, error) {
	return query.Period{
		SobId:        po.SobId,
		Id:           po.Id,
		FiscalYear:   po.FiscalYear,
		PeriodNumber: po.PeriodNumber,
		OpeningTime:  po.OpeningTime,
		EndingTime:   po.EndingTime,
		IsClosed:     po.IsClosed,
		CreatedAt:    po.CreatedAt,
		UpdatedAt:    po.UpdatedAt,
	}, nil
}

func ledgerBOToPO(bo ledger.Ledger) ledgerPO {
	return ledgerPO{
		Id:             bo.Id(),
		SobId:          bo.SobId(),
		AccountId:      bo.AccountId(),
		PeriodId:       bo.PeriodId(),
		OpeningBalance: bo.OpeningBalance(),
		EndingBalance:  bo.EndingBalance(),
		PeriodDebit:    bo.PeriodDebit(),
		PeriodCredit:   bo.PeriodCredit(),
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
		po.AccountId,
		po.PeriodId,
		po.OpeningBalance,
		po.EndingBalance,
		po.PeriodDebit,
		po.PeriodCredit,
		*accountBO,
	)
}

func ledgerPOToDTO(po ledgerPO) (query.Ledger, error) {
	accountDTO, err := accountPOToDTO(po.Account)
	if err != nil {
		return query.Ledger{}, err
	}

	return query.Ledger{
		Id:             po.Id,
		SobId:          po.SobId,
		AccountId:      po.AccountId,
		PeriodId:       po.PeriodId,
		OpeningBalance: po.OpeningBalance,
		EndingBalance:  po.EndingBalance,
		PeriodDebit:    po.PeriodDebit,
		PeriodCredit:   po.PeriodCredit,
		Account:        accountDTO,
		CreatedAt:      po.CreatedAt,
		UpdatedAt:      po.UpdatedAt,
	}, nil
}
