package db

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"github/fims-proto/fims-proto-ms/internal/sob/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SobPostgresRepository struct{}

func NewSobPostgresRepository() *SobPostgresRepository {
	return &SobPostgresRepository{}
}

func (r SobPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	if err := db.AutoMigrate(&sob{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r SobPostgresRepository) CreateSob(ctx context.Context, sob *domain.Sob) error {
	db := readDBFromCtx(ctx)

	dbSob, err := marshall(sob)
	if err != nil {
		return errors.Wrap(err, "failed to marshal sob")
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(dbSob).Error
	}); err != nil {
		return errors.Wrap(err, "create sob failed")
	}

	return nil
}

func (r SobPostgresRepository) UpdateSob(ctx context.Context, sobId uuid.UUID, updateFn func(s *domain.Sob) (*domain.Sob, error)) error {
	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		dbSob := &sob{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(dbSob, "id = ?", sobId).Error; err != nil {
			return errors.Wrap(err, "failed to find sob")
		}

		domainSob, err := unmarshallToDomain(dbSob)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshall sob")
		}

		updatedDomainSob, err := updateFn(domainSob)
		if err != nil {
			return errors.Wrap(err, "failed to update sob in transaction")
		}

		dbSob, err = marshall(updatedDomainSob)
		if err != nil {
			return errors.Wrap(err, "failed to marshal sob")
		}
		if err := tx.Save(dbSob).Error; err != nil {
			return errors.Wrap(err, "failed to save sob")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "failed to updated sob")
	}
	return nil
}

func (r SobPostgresRepository) AllSobs(ctx context.Context) ([]query.Sob, error) {
	db := readDBFromCtx(ctx)

	var dbSobs []sob
	if err := db.Find(&dbSobs).Error; err != nil {
		return []query.Sob{}, errors.Wrap(err, "failed to read all sob")
	}

	var querySobs []query.Sob
	for _, dbSob := range dbSobs {
		querySob, err := unmarshallToQuery(&dbSob)
		if err != nil {
			return []query.Sob{}, errors.Wrap(err, "failed to unmarshall sob")
		}
		querySobs = append(querySobs, querySob)
	}
	return querySobs, nil
}

func (r SobPostgresRepository) SobById(ctx context.Context, sobId uuid.UUID) (query.Sob, error) {
	db := readDBFromCtx(ctx)

	dbSob := sob{}
	if err := db.Where("id = ?", sobId).First(&dbSob).Error; err != nil {
		return query.Sob{}, errors.Wrapf(err, "failed to read sob %s", sobId)
	}

	querySob, err := unmarshallToQuery(&dbSob)
	if err != nil {
		return query.Sob{}, errors.Wrap(err, "failed to unmarshall sob")
	}
	return querySob, nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
