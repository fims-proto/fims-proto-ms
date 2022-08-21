package db

import (
	"encoding/json"
	"time"

	"github/fims-proto/fims-proto-ms/internal/user/domain/user"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/user/app/query"
)

type userPO struct {
	Id        uuid.UUID `gorm:"type:uuid"`
	Traits    pgtype.JSONB
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// table names

func (u userPO) TableName() string {
	return "a_users"
}

// mappers

func userBOToPO(bo user.User) (userPO, error) {
	var traits pgtype.JSONB
	if err := traits.Set(bo.Traits()); err != nil {
		return userPO{}, errors.Wrap(err, "convert json.RawMessage to pgtype.JSONB failed")
	}

	return userPO{
		Id:     bo.Id(),
		Traits: traits,
	}, nil
}

func userPOToBO(po userPO) (*user.User, error) {
	var traits json.RawMessage
	marshalJSON, err := po.Traits.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "convert pgtype.JSONB to json.RawMessage failed")
	}

	if err = traits.UnmarshalJSON(marshalJSON); err != nil {
		return nil, errors.Wrap(err, "convert pgtype.JSONB to json.RawMessage failed")
	}

	return user.New(po.Id, traits)
}

func userPOToDTO(po userPO) (query.User, error) {
	var traits json.RawMessage
	marshalJSON, err := po.Traits.MarshalJSON()
	if err != nil {
		return query.User{}, errors.Wrap(err, "convert pgtype.JSONB to json.RawMessage failed")
	}

	if err = traits.UnmarshalJSON(marshalJSON); err != nil {
		return query.User{}, errors.Wrap(err, "convert pgtype.JSONB to json.RawMessage failed")
	}

	return query.User{
		Id:        po.Id,
		Traits:    traits,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}, nil
}
