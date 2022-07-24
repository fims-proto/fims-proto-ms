package db

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/user/app/query"
	"github/fims-proto/fims-proto-ms/internal/user/domain"
)

type user struct {
	Id        uuid.UUID `gorm:"type:uuid"`
	Traits    pgtype.JSONB
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func marshal(u domain.User) (user, error) {
	var traits pgtype.JSONB
	if err := traits.Set(u.Traits()); err != nil {
		return user{}, errors.Wrap(err, "convert json.RawMessage to pgtype.JSONB failed")
	}
	return user{
		Id:     u.Id(),
		Traits: traits,
	}, nil
}

func unmarshalToDomain(dbu user) (*domain.User, error) {
	var traits json.RawMessage
	marshalJSON, err := dbu.Traits.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(err, "convert pgtype.JSONB to json.RawMessage failed")
	}
	if err = traits.UnmarshalJSON(marshalJSON); err != nil {
		return nil, errors.Wrap(err, "convert pgtype.JSONB to json.RawMessage failed")
	}
	return domain.NewUser(dbu.Id, traits)
}

func unmarshalToQuery(dbu user) (query.User, error) {
	var traits json.RawMessage
	marshalJSON, err := dbu.Traits.MarshalJSON()
	if err != nil {
		return query.User{}, errors.Wrap(err, "convert pgtype.JSONB to json.RawMessage failed")
	}
	if err = traits.UnmarshalJSON(marshalJSON); err != nil {
		return query.User{}, errors.Wrap(err, "convert pgtype.JSONB to json.RawMessage failed")
	}
	return query.User{
		Id:        dbu.Id,
		Traits:    traits,
		CreatedAt: dbu.CreatedAt,
		UpdatedAt: dbu.UpdatedAt,
	}, nil
}
