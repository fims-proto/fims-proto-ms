package db

import (
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
	"time"

	"github.com/google/uuid"
)

type account struct {
	Id             uuid.UUID `gorm:"type:uuid"`
	SobId          string    `gorm:"uniqueIndex:accounts_sobid_number_key"`
	Number         string    `gorm:"uniqueIndex:accounts_sobid_number_key"`
	Title          string
	SuperiorNumber string
	AccountType    string
	CreatedAt      time.Time `gorm:"<-:create"`
	UpdatedAt      time.Time
}

func marshall(a *domain.Account) *account {
	return &account{
		Id:             a.Id(),
		SobId:          a.Sob(),
		Number:         a.Number(),
		Title:          a.Title(),
		SuperiorNumber: a.SuperiorNumber(),
		AccountType:    a.Type().String(),
	}
}

func unmarshallToQuery(dba *account) query.Account {
	return query.Account{
		Sob:             dba.SobId,
		Number:          dba.Number,
		Title:           dba.Title,
		AccountType:     dba.AccountType,
		SuperiorAccount: nil,
	}
}
