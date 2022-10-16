package db

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/sob/domain/sob"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
)

type sobPO struct {
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

// table names

func (s sobPO) TableName() string {
	return "a_sobs"
}

// schemas

func (s sobPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return s.TableName(), nil
	}
	return "", errors.Errorf("sobPO doesn't have association named %s", entity)
}

// mappers

func sobBOToPO(bo sob.Sob) (sobPO, error) {
	var intArray pgtype.Int4Array
	if err := intArray.Set(bo.AccountsCodeLength()); err != nil {
		return sobPO{}, errors.Wrap(err, "convert []int to Int4Array failed")
	}

	return sobPO{
		Id:                  bo.Id(),
		Name:                bo.Name(),
		Description:         bo.Description(),
		BaseCurrency:        bo.BaseCurrency(),
		StartingPeriodYear:  bo.StartingPeriodYear(),
		StartingPeriodMonth: bo.StartingPeriodMonth(),
		AccountsCodeLength:  intArray,
	}, nil
}

func sobPOToBO(po sobPO) (*sob.Sob, error) {
	var codesLength []int
	if err := po.AccountsCodeLength.AssignTo(&codesLength); err != nil {
		return nil, errors.Wrap(err, "assign Int4Array to []int failed")
	}

	return sob.New(
		po.Id,
		po.Name,
		po.Description,
		po.BaseCurrency,
		po.StartingPeriodYear,
		po.StartingPeriodMonth,
		codesLength, // from 4-2-2 to [4,2,2]
	)
}

func sobPOToDTO(po sobPO) (query.Sob, error) {
	var codesLength []int
	if err := po.AccountsCodeLength.AssignTo(&codesLength); err != nil {
		return query.Sob{}, errors.Wrap(err, "assign Int4Array to []int failed")
	}

	return query.Sob{
		Id:                  po.Id,
		Name:                po.Name,
		Description:         po.Description,
		BaseCurrency:        po.BaseCurrency,
		StartingPeriodYear:  po.StartingPeriodYear,
		StartingPeriodMonth: po.StartingPeriodMonth,
		AccountsCodeLength:  codesLength, // from 4-2-2 to [4,2,2]
	}, nil
}
