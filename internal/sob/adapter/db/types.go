package db

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"github/fims-proto/fims-proto-ms/internal/sob/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

type sob struct {
	Id                  uuid.UUID `gorm:"type:uuid"`
	Name                string    `gorm:"uniqueIndex"`
	Description         string
	BaseCurrency        string
	StartingPeriodYear  int
	StartingPeriodMonth int
	AccountsCodeLength  pgtype.Int4Array `gorm:"type:integer[]"`
	CreatedAt           time.Time        `gorm:"<-:create"`
	UpdatedAt           time.Time
}

func marshal(s domain.Sob) (sob, error) {
	var intArray pgtype.Int4Array
	if err := intArray.Set(s.AccountsCodeLength()); err != nil {
		return sob{}, errors.Wrap(err, "convert []int to Int4Array failed")
	}
	return sob{
		Id:                  s.Id(),
		Name:                s.Name(),
		Description:         s.Description(),
		BaseCurrency:        s.BaseCurrency(),
		StartingPeriodYear:  s.StartingPeriodYear(),
		StartingPeriodMonth: s.StartingPeriodMonth(),
		AccountsCodeLength:  intArray,
	}, nil
}

func unmarshalToDomain(dbs sob) (*domain.Sob, error) {
	var codesLength []int
	if err := dbs.AccountsCodeLength.AssignTo(&codesLength); err != nil {
		return nil, errors.Wrap(err, "assign Int4Array to []int failed")
	}
	return domain.NewSob(
		dbs.Id,
		dbs.Name,
		dbs.Description,
		dbs.BaseCurrency,
		dbs.StartingPeriodYear,
		dbs.StartingPeriodMonth,
		codesLength, // from 4-2-2 to [4,2,2]
	)
}

func unmarshalToQuery(dbs sob) (query.Sob, error) {
	var codesLength []int
	if err := dbs.AccountsCodeLength.AssignTo(&codesLength); err != nil {
		return query.Sob{}, errors.Wrap(err, "assign Int4Array to []int failed")
	}
	return query.Sob{
		Id:                  dbs.Id,
		Name:                dbs.Name,
		Description:         dbs.Description,
		BaseCurrency:        dbs.BaseCurrency,
		StartingPeriodYear:  dbs.StartingPeriodYear,
		StartingPeriodMonth: dbs.StartingPeriodMonth,
		AccountsCodeLength:  codesLength, // from 4-2-2 to [4,2,2]
	}, nil
}
