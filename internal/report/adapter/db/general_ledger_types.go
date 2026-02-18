package db

import (
	"github/fims-proto/fims-proto-ms/internal/report/domain/general_ledger"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

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
}

type accountPO struct {
	Id                uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId             uuid.UUID `gorm:"type:uuid"`
	SuperiorAccountId uuid.UUID `gorm:"type:uuid"`
	Title             string
	AccountNumber     string
	Level             int
	IsLeaf            bool
	Class             int
	Group             int
	BalanceDirection  string
}

type periodPO struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId        uuid.UUID `gorm:"type:uuid"`
	FiscalYear   int
	PeriodNumber int
}

// table name

func (l ledgerPO) TableName() string {
	return "a_ledgers"
}

func (a accountPO) TableName() string {
	return "a_accounts"
}

func (a periodPO) TableName() string {
	return "a_periods"
}

// mappers

func ledgerPOToBO(po ledgerPO) (*general_ledger.Ledger, error) {
	account, err := accountPOToBO(po.Account)
	if err != nil {
		return nil, err
	}

	return general_ledger.NewLedger(
		account,
		periodPOToBO(po.Period),
		po.OpeningDebitBalance,
		po.OpeningCreditBalance,
		po.PeriodDebit,
		po.PeriodCredit,
		po.EndingDebitBalance,
		po.EndingCreditBalance,
	), nil
}

func accountPOToBO(po accountPO) (*general_ledger.Account, error) {
	return general_ledger.NewAccount(po.BalanceDirection)
}

func periodPOToBO(po periodPO) *general_ledger.Period {
	return general_ledger.NewPeriod(
		po.FiscalYear,
		po.PeriodNumber,
	)
}
