package db

import (
	"encoding/json"
	"time"

	"github/fims-proto/fims-proto-ms/internal/user/domain/user"

	"github/fims-proto/fims-proto-ms/internal/user/app/query"

	"github.com/google/uuid"
)

type userPO struct {
	Id     uuid.UUID       `gorm:"type:uuid"`
	Traits json.RawMessage `gorm:"type:jsonb"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// mappers

func userBOToPO(bo user.User) userPO {
	return userPO{
		Id:     bo.Id(),
		Traits: bo.Traits(),
	}
}

func userPOToBO(po userPO) (*user.User, error) {
	return user.New(po.Id, po.Traits)
}

func userPOToDTO(po userPO) query.User {
	return query.User{
		Id:        po.Id,
		Traits:    po.Traits,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}
