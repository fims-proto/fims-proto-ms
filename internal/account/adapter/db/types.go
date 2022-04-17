package db

import (
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
	"time"

	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

type account struct {
	Id                uuid.UUID `gorm:"type:uuid"`
	SobId             uuid.UUID `gorm:"type:uuid;uniqueIndex:accounts_sobid_sai_number_key"`
	SuperiorAccountId uuid.UUID `gorm:"type:uuid;uniqueIndex:accounts_sobid_sai_number_key"`
	LevelNumber       int       `gorm:"uniqueIndex:accounts_sobid_sai_number_key"`
	Title             string
	Level             int
	AccountType       string
	SuperiorNumbers   pgtype.Int4Array `gorm:"type:integer[]"`
	BalanceDirection  string
	CreatedAt         time.Time `gorm:"<-:create"`
	UpdatedAt         time.Time
}

func marshall(a *domain.Account) (*account, error) {
	var int4array pgtype.Int4Array
	if err := int4array.Set(a.SuperiorNumbers()); err != nil {
		return nil, errors.Wrap(err, "convert []int to Int4Array failed")
	}
	return &account{
		Id:                a.Id(),
		SobId:             a.SobId(),
		SuperiorAccountId: a.SuperiorAccountId(),
		LevelNumber:       a.LevelNumber(),
		Title:             a.Title(),
		Level:             a.Level(),
		AccountType:       a.Type().String(),
		SuperiorNumbers:   int4array,
		BalanceDirection:  a.BalanceDirection().String(),
	}, nil
}

func unmarshallToQuery(dba *account) (*query.Account, error) {
	var numbers []int
	if err := dba.SuperiorNumbers.AssignTo(&numbers); err != nil {
		return nil, errors.Wrap(err, "assign Int4Array to []int failed")
	}
	accountType, err := commonAccount.NewAccountType(dba.AccountType)
	if err != nil {
		return nil, errors.Wrap(err, "should not happen: failed to parse account type")
	}
	direction, err := commonAccount.NewDirection(dba.BalanceDirection)
	if err != nil {
		return nil, errors.Wrap(err, "should not happen: failed to parse balance direction")
	}
	return &query.Account{
		Id:                dba.Id,
		SobId:             dba.SobId,
		SuperiorAccountId: dba.SuperiorAccountId,
		SuperiorNumbers:   numbers,
		LevelNumber:       dba.LevelNumber,
		Title:             dba.Title,
		Level:             dba.Level,
		AccountType:       accountType,
		SuperiorAccount:   nil,
		BalanceDirection:  direction,
	}, nil
}
