package db

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/database"

	"github/fims-proto/fims-proto-ms/internal/sob/domain/sob"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SobPostgresRepository struct{}

func NewSobPostgresRepository() *SobPostgresRepository {
	return &SobPostgresRepository{}
}

func (r SobPostgresRepository) Migrate(ctx context.Context) error {
	db := database.ReadDBFromContext(ctx)

	if err := db.AutoMigrate(&sobPO{}); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	return nil
}

func (r SobPostgresRepository) CreateSob(ctx context.Context, sob *sob.Sob) error {
	db := database.ReadDBFromContext(ctx)

	po, err := sobBOToPO(*sob)
	if err != nil {
		return err
	}

	return db.Create(&po).Error
}

func (r SobPostgresRepository) UpdateSob(ctx context.Context, sobId uuid.UUID, updateFn func(s *sob.Sob) (*sob.Sob, error)) error {
	db := database.ReadDBFromContext(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		po := sobPO{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po, "id = ?", sobId).Error; err != nil {
			return fmt.Errorf("failed to find sob: %w", err)
		}

		bo, err := sobPOToBO(po)
		if err != nil {
			return err
		}

		updatedBO, err := updateFn(bo)
		if err != nil {
			return fmt.Errorf("failed to update sob: %w", err)
		}

		po, err = sobBOToPO(*updatedBO)
		if err != nil {
			return err
		}

		return tx.Save(&po).Error
	})
}

// Queries

func (r SobPostgresRepository) SearchSobs(ctx context.Context, pageRequest data.PageRequest) (data.Page[query.Sob], error) {
	return data.SearchEntities(ctx, pageRequest, sobPO{}, sobPOToDTO, database.ReadDBFromContext(ctx))
}

func (r SobPostgresRepository) SobById(ctx context.Context, sobId uuid.UUID) (query.Sob, error) {
	db := database.ReadDBFromContext(ctx)

	dbSob := sobPO{}
	if err := db.Where("id = ?", sobId).First(&dbSob).Error; err != nil {
		return query.Sob{}, fmt.Errorf("failed to read sob %s: %w", sobId, err)
	}

	return sobPOToDTO(dbSob), nil
}
