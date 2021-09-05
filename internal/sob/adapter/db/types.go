package db

import (
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"github/fims-proto/fims-proto-ms/internal/sob/domain"
	"time"
)

type sob struct {
	Id          string
	Name        string `gorm:"uniqueIndex"`
	Description string
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
}

func marshall(s *domain.Sob) *sob {
	return &sob{
		Id:          s.Id(),
		Name:        s.Name(),
		Description: s.Description(),
	}
}

func unmarshallToDomain(dbs *sob) (*domain.Sob, error) {
	return domain.NewSob(dbs.Id, dbs.Name, dbs.Description)
}

func unmarshallToQuery(dbs *sob) query.Sob {
	return query.Sob{
		Id:          dbs.Id,
		Name:        dbs.Name,
		Description: dbs.Description,
	}
}
