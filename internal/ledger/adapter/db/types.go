package db

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type period struct {
	Id               uuid.UUID `gorm:"type:uuid"`
	SobId            uuid.UUID `gorm:"type:uuid;uniqueIndex:periods_sobid_year_number_key"`
	PreviousPeriodId uuid.UUID `gorm:"type:uuid"`
	FinancialYear    int       `gorm:"uniqueIndex:periods_sobid_year_number_key"`
	Number           int       `gorm:"uniqueIndex:periods_sobid_year_number_key"`
	OpeningTime      time.Time
	EndingTime       time.Time
	IsClosed         bool
	CreatedAt        time.Time `gorm:"<-:create"`
	UpdatedAt        time.Time
}

type ledger struct {
	Id             uuid.UUID `gorm:"type:uuid"`
	PeriodId       uuid.UUID `gorm:"type:uuid;uniqueIndex:ledgers_periodid_accountid_key"`
	AccountId      uuid.UUID `gorm:"type:uuid;uniqueIndex:ledgers_periodid_accountid_key"`
	OpeningBalance decimal.Decimal
	EndingBalance  decimal.Decimal
	Debit          decimal.Decimal
	Credit         decimal.Decimal
	CreatedAt      time.Time `gorm:"<-:create"`
	UpdatedAt      time.Time
}

type ledgerLog struct {
	Id              uuid.UUID `gorm:"type:uuid"`
	VoucherId       uuid.UUID `gorm:"type:uuid"`
	PostingId       uuid.UUID `gorm:"type:uuid"`
	AccountId       uuid.UUID `gorm:"type:uuid"`
	TransactionTime time.Time
	Debit           decimal.Decimal
	Credit          decimal.Decimal
	CreatedAt       time.Time `gorm:"<-:create"`
	UpdatedAt       time.Time
}

func marshallPeriod(p *domain.Period) *period {
	return &period{
		Id:               p.Id(),
		SobId:            p.SobId(),
		PreviousPeriodId: p.PreviousPeriodId(),
		FinancialYear:    p.FinancialYear(),
		Number:           p.Number(),
		OpeningTime:      p.OpeningTime(),
		EndingTime:       p.EndingTime(),
		IsClosed:         p.IsClosed(),
	}
}

func marshallLedger(l *domain.Ledger) *ledger {
	return &ledger{
		Id:             l.Id(),
		PeriodId:       l.PeriodId(),
		AccountId:      l.AccountId(),
		OpeningBalance: l.OpeningBalance(),
		EndingBalance:  l.EndingBalance(),
		Debit:          l.Debit(),
		Credit:         l.Credit(),
	}
}

func marshallLedgerLog(l *domain.LedgerLog) *ledgerLog {
	return &ledgerLog{
		Id:              l.Id(),
		VoucherId:       l.VoucherId(),
		PostingId:       l.PostingId(),
		AccountId:       l.AccountId(),
		TransactionTime: l.TransactionTime(),
		Debit:           l.Debit(),
		Credit:          l.Credit(),
	}
}

func unmarshallPeriodToDomain(dbp *period) (*domain.Period, error) {
	return domain.NewPeriod(
		dbp.Id,
		dbp.SobId,
		dbp.PreviousPeriodId,
		dbp.FinancialYear,
		dbp.Number,
		dbp.OpeningTime,
		dbp.EndingTime,
		dbp.IsClosed,
	)
}

func unmarshallLedgerToDomain(dbl *ledger) (*domain.Ledger, error) {
	return domain.NewLedger(
		dbl.Id,
		dbl.PeriodId,
		dbl.AccountId,
		dbl.OpeningBalance,
		dbl.EndingBalance,
		dbl.Debit,
		dbl.Credit,
	)
}

func unmarshallPeriodToQuery(dbp *period) query.Period {
	return query.Period{
		Id:               dbp.Id,
		SobId:            dbp.SobId,
		PreviousPeriodId: dbp.PreviousPeriodId,
		FinancialYear:    dbp.FinancialYear,
		Number:           dbp.Number,
		OpeningTime:      dbp.OpeningTime,
		EndingTime:       dbp.EndingTime,
		IsClosed:         dbp.IsClosed,
		CreatedAt:        dbp.CreatedAt,
		UpdatedAt:        dbp.UpdatedAt,
	}
}

func unmarshallLedgerToQuery(dbl *ledger) query.Ledger {
	return query.Ledger{
		Id:             dbl.Id,
		PeriodId:       dbl.PeriodId,
		Account:        query.Account{Id: dbl.AccountId},
		OpeningBalance: dbl.OpeningBalance,
		EndingBalance:  dbl.EndingBalance,
		Debit:          dbl.Debit,
		Credit:         dbl.Credit,
		CreatedAt:      dbl.CreatedAt,
		UpdatedAt:      dbl.UpdatedAt,
	}
}

func unmarshallLedgerLogToQuery(dbl *ledgerLog) query.LedgerLog {
	return query.LedgerLog{
		Id:              dbl.Id,
		VoucherId:       dbl.VoucherId,
		PostingId:       dbl.PostingId,
		AccountId:       dbl.AccountId,
		TransactionTime: dbl.TransactionTime,
		Debit:           dbl.Debit,
		Credit:          dbl.Credit,
		CreatedAt:       dbl.CreatedAt,
		UpdatedAt:       dbl.UpdatedAt,
	}
}
