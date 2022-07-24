package db

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain"

	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

type account struct {
	Id                uuid.UUID        `gorm:"type:uuid"`
	SobId             uuid.UUID        `gorm:"type:uuid;uniqueIndex:accounts_sobid_number_key"`
	AccountNumber     string           `gorm:"uniqueIndex:accounts_sobid_number_key"`
	SuperiorAccountId uuid.UUID        `gorm:"type:uuid"`
	NumberHierarchy   pgtype.Int4Array `gorm:"type:integer[]"`
	Title             string
	Level             int
	AccountType       string
	BalanceDirection  string
	CreatedAt         time.Time `gorm:"<-:create"`
	UpdatedAt         time.Time
}

func marshal(a domain.Account) (account, error) {
	var int4array pgtype.Int4Array
	if err := int4array.Set(a.NumberHierarchy()); err != nil {
		return account{}, errors.Wrap(err, "convert []int to Int4Array failed")
	}
	return account{
		Id:                a.Id(),
		SobId:             a.SobId(),
		AccountNumber:     a.AccountNumber(),
		SuperiorAccountId: a.SuperiorAccountId(),
		NumberHierarchy:   int4array,
		Title:             a.Title(),
		Level:             a.Level(),
		AccountType:       a.Type().String(),
		BalanceDirection:  a.BalanceDirection().String(),
	}, nil
}

func unmarshalToQuery(dba account) (query.Account, error) {
	var numbers []int
	if err := dba.NumberHierarchy.AssignTo(&numbers); err != nil {
		return query.Account{}, errors.Wrap(err, "assign Int4Array to []int failed")
	}
	accountType, err := commonAccount.NewAccountType(dba.AccountType)
	if err != nil {
		return query.Account{}, errors.Wrap(err, "should not happen: failed to parse account type")
	}
	direction, err := commonAccount.NewDirection(dba.BalanceDirection)
	if err != nil {
		return query.Account{}, errors.Wrap(err, "should not happen: failed to parse balance direction")
	}
	res := query.Account{
		Id:               dba.Id,
		SobId:            dba.SobId,
		Title:            dba.Title,
		AccountNumber:    dba.AccountNumber,
		NumberHierarchy:  numbers,
		AccountType:      accountType,
		BalanceDirection: direction,
		Level:            dba.Level,
		SuperiorAccount:  nil,
		CreatedAt:        dba.CreatedAt,
		UpdatedAt:        dba.UpdatedAt,
	}
	if dba.SuperiorAccountId != uuid.Nil {
		res.SuperiorAccount = &query.Account{Id: dba.SuperiorAccountId}
	}
	return res, nil
}
