package db

import (
	"time"

	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account_configuration"
	"github/fims-proto/fims-proto-ms/internal/account/domain/period"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

type accountConfigurationPO struct {
	SobId             uuid.UUID `gorm:"type:uuid;uniqueIndex:accountConfigurations_sobid_number_key"`
	AccountId         uuid.UUID `gorm:"type:uuid;primaryKey"`
	SuperiorAccountId uuid.UUID `gorm:"type:uuid"`
	Title             string
	AccountNumber     string           `gorm:"uniqueIndex:accountConfigurations_sobid_number_key"`
	NumberHierarchy   pgtype.Int4Array `gorm:"type:integer[]"`
	Level             int
	AccountType       string
	BalanceDirection  string
	CreatedAt         time.Time `gorm:"<-:create"`
	UpdatedAt         time.Time
}

type periodPO struct {
	SobId            uuid.UUID `gorm:"type:uuid;uniqueIndex:periods_sobid_year_number_key"`
	PeriodId         uuid.UUID `gorm:"type:uuid;primaryKey"`
	PreviousPeriodId uuid.UUID `gorm:"type:uuid"`
	FinancialYear    int       `gorm:"uniqueIndex:periods_sobid_year_number_key"`
	Number           int       `gorm:"uniqueIndex:periods_sobid_year_number_key"`
	OpeningTime      time.Time
	EndingTime       time.Time
	IsClosed         bool
	CreatedAt        time.Time `gorm:"<-:create"`
	UpdatedAt        time.Time
}

type accountPO struct {
	SobId          uuid.UUID `gorm:"type:uuid"`
	AccountId      uuid.UUID `gorm:"type:uuid;primaryKey"`
	PeriodId       uuid.UUID `gorm:"type:uuid;primaryKey"`
	OpeningBalance decimal.Decimal
	EndingBalance  decimal.Decimal
	PeriodDebit    decimal.Decimal
	PeriodCredit   decimal.Decimal
	CreatedAt      time.Time `gorm:"<-:create"`
	UpdatedAt      time.Time
}

type accountVO struct {
	Account accountPO              `gorm:"embedded"`
	Config  accountConfigurationPO `gorm:"embedded"`
	Period  periodPO               `gorm:"embedded"`
}

// table names

func (a accountPO) TableName() string {
	return "accounts"
}

func (a accountConfigurationPO) TableName() string {
	return "account_configurations"
}

func (p periodPO) TableName() string {
	return "periods"
}

// mappers

func accountConfigurationBOToPO(bo account_configuration.AccountConfiguration) (accountConfigurationPO, error) {
	var int4array pgtype.Int4Array
	if err := int4array.Set(bo.NumberHierarchy()); err != nil {
		return accountConfigurationPO{}, errors.Wrap(err, "convert []int to Int4Array failed")
	}
	return accountConfigurationPO{
		SobId:             bo.SobId(),
		AccountId:         bo.AccountId(),
		SuperiorAccountId: bo.SuperiorAccountId(),
		Title:             bo.Title(),
		AccountNumber:     bo.AccountNumber(),
		NumberHierarchy:   int4array,
		Level:             bo.Level(),
		AccountType:       bo.AccountType().String(),
		BalanceDirection:  bo.BalanceDirection().String(),
	}, nil
}

func accountConfigurationPOToBO(po accountConfigurationPO) (*account_configuration.AccountConfiguration, error) {
	var numberHierarchy []int
	if err := po.NumberHierarchy.AssignTo(&numberHierarchy); err != nil {
		return nil, errors.Wrap(err, "assign Int4Array to []int failed")
	}

	return account_configuration.New(
		po.SobId,
		po.AccountId,
		po.SuperiorAccountId,
		po.Title,
		po.AccountNumber,
		numberHierarchy,
		po.Level,
		po.AccountType,
		po.BalanceDirection,
	)
}

func accountConfigurationPOToDTO(po accountConfigurationPO) (query.AccountConfiguration, error) {
	var numberHierarchy []int
	if err := po.NumberHierarchy.AssignTo(&numberHierarchy); err != nil {
		return query.AccountConfiguration{}, errors.Wrap(err, "assign Int4Array to []int failed")
	}

	return query.AccountConfiguration{
		SobId:             po.SobId,
		AccountId:         po.AccountId,
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
		SobId:            bo.SobId(),
		PeriodId:         bo.PeriodId(),
		PreviousPeriodId: bo.PreviousPeriodId(),
		FinancialYear:    bo.FinancialYear(),
		Number:           bo.Number(),
		OpeningTime:      bo.OpeningTime(),
		EndingTime:       bo.EndingTime(),
		IsClosed:         bo.IsClosed(),
	}
}

func periodPOToDTO(po periodPO) query.Period {
	return query.Period{
		SobId:            po.SobId,
		PeriodId:         po.PeriodId,
		PreviousPeriodId: po.PreviousPeriodId,
		FinancialYear:    po.FinancialYear,
		Number:           po.Number,
		OpeningTime:      po.OpeningTime,
		EndingTime:       po.EndingTime,
		IsClosed:         po.IsClosed,
		CreatedAt:        po.CreatedAt,
		UpdatedAt:        po.UpdatedAt,
	}
}

func accountBOToPO(bo account.Account) accountPO {
	return accountPO{
		SobId:          bo.SobId(),
		AccountId:      bo.AccountId(),
		PeriodId:       bo.PeriodId(),
		OpeningBalance: bo.OpeningBalance(),
		EndingBalance:  bo.EndingBalance(),
		PeriodDebit:    bo.PeriodDebit(),
		PeriodCredit:   bo.PeriodCredit(),
	}
}

func accountVOToBO(vo accountVO) (*account.Account, error) {
	configuration, err := accountConfigurationPOToBO(vo.Config)
	if err != nil {
		return nil, err
	}

	return account.New(
		vo.Account.SobId,
		vo.Account.AccountId,
		vo.Account.PeriodId,
		vo.Account.OpeningBalance,
		vo.Account.EndingBalance,
		vo.Account.PeriodDebit,
		vo.Account.PeriodCredit,
		*configuration,
	)
}

func accountVOToDTO(vo accountVO) (query.Account, error) {
	configuration, err := accountConfigurationPOToDTO(vo.Config)
	if err != nil {
		return query.Account{}, err
	}

	return query.Account{
		SobId:          vo.Account.SobId,
		AccountId:      vo.Account.AccountId,
		OpeningBalance: vo.Account.OpeningBalance,
		EndingBalance:  vo.Account.EndingBalance,
		PeriodDebit:    vo.Account.PeriodDebit,
		PeriodCredit:   vo.Account.PeriodCredit,
		CreatedAt:      vo.Account.CreatedAt,
		UpdatedAt:      vo.Account.UpdatedAt,
		Configuration:  configuration,
		Period:         periodPOToDTO(vo.Period),
	}, nil
}
