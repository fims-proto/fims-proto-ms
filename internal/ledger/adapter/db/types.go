package db

import (
	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ledger struct {
	Id             uuid.UUID `gorm:"type:uuid"`
	SobId          string    `gorm:"uniqueIndex:ledgers_sobid_number_key"`
	Number         string    `gorm:"uniqueIndex:ledgers_sobid_number_key"`
	Title          string
	SuperiorNumber string
	AccountType    string
	Debit          decimal.Decimal
	Credit         decimal.Decimal
	Balance        decimal.Decimal
	CreatedAt      time.Time `gorm:"<-:create"`
	UpdatedAt      time.Time
}

func marshall(l *domain.Ledger) *ledger {
	return &ledger{
		Id:             l.Id(),
		SobId:          l.Sob(),
		Number:         l.Number(),
		Title:          l.Title(),
		SuperiorNumber: l.SuperiorNumber(),
		AccountType:    l.AccountType().String(),
		Debit:          l.Debit(),
		Credit:         l.Credit(),
		Balance:        l.Balance(),
	}
}

func unmarshallToDomain(dbl *ledger) (*domain.Ledger, error) {
	return domain.NewLedger(dbl.Id, dbl.SobId, dbl.Number, dbl.Title, dbl.SuperiorNumber, dbl.AccountType, dbl.Debit, dbl.Credit, dbl.Balance)
}

func unmarshallToQuery(dbl *ledger) query.Ledger {
	return query.Ledger{
		Id:             dbl.Id,
		Sob:            dbl.SobId,
		Number:         dbl.Number,
		Title:          dbl.Title,
		SuperiorNumber: dbl.SuperiorNumber,
		AccountType:    dbl.AccountType,
		Debit:          dbl.Debit,
		Credit:         dbl.Credit,
		Balance:        dbl.Balance,
	}
}
