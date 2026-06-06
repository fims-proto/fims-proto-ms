package db

import (
	"fmt"
	"time"

	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"github/fims-proto/fims-proto-ms/internal/sob/domain/sob"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type sobPO struct {
	Id                  uuid.UUID `gorm:"type:uuid"`
	Name                string    `gorm:"uniqueIndex:UQ_Sobs_Name"`
	Description         string
	BaseCurrency        string
	StartingPeriodYear  int
	StartingPeriodMonth int
	AccountsCodeLength  pgtype.Int4Array `gorm:"type:integer[]"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// schemas

func (s sobPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return "sobs", nil
	}
	return "", fmt.Errorf("sobPO doesn't have association named %s", entity)
}

// mappers

func sobBOToPO(bo sob.Sob) sobPO {
	var intArray pgtype.Int4Array
	if err := intArray.Set(bo.AccountsCodeLength()); err != nil {
		panic(fmt.Errorf("failed to convert []int to Int4Array: %w", err))
	}

	return sobPO{
		Id:                  bo.Id(),
		Name:                bo.Name(),
		Description:         bo.Description(),
		BaseCurrency:        bo.BaseCurrency(),
		StartingPeriodYear:  bo.StartingPeriodYear(),
		StartingPeriodMonth: bo.StartingPeriodMonth(),
		AccountsCodeLength:  intArray,
	}
}

func sobPOToBO(po sobPO) (*sob.Sob, error) {
	var codesLength []int
	if err := po.AccountsCodeLength.AssignTo(&codesLength); err != nil {
		return nil, fmt.Errorf("failed to assign Int4Array to []int: %w", err)
	}

	return sob.New(
		po.Id,
		po.Name,
		po.Description,
		po.BaseCurrency,
		po.StartingPeriodYear,
		po.StartingPeriodMonth,
		codesLength,
	)
}

func sobPOToDTO(po sobPO) query.Sob {
	var codesLength []int
	if err := po.AccountsCodeLength.AssignTo(&codesLength); err != nil {
		panic(fmt.Errorf("failed to assign Int4Array to []int: %w", err))
	}

	return query.Sob{
		Id:                  po.Id,
		Name:                po.Name,
		Description:         po.Description,
		BaseCurrency:        po.BaseCurrency,
		StartingPeriodYear:  po.StartingPeriodYear,
		StartingPeriodMonth: po.StartingPeriodMonth,
		AccountsCodeLength:  codesLength,
		CreatedAt:           po.CreatedAt,
		UpdatedAt:           po.UpdatedAt,
	}
}
