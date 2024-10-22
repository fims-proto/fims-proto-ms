package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"github/fims-proto/fims-proto-ms/internal/sob/domain/sob"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SobPostgresRepository struct {
	dataSource datasource.DataSource
}

func NewSobPostgresRepository(dataSource datasource.DataSource) *SobPostgresRepository {
	return &SobPostgresRepository{
		dataSource: dataSource,
	}
}

func (r SobPostgresRepository) Migrate(ctx context.Context) error {
	db := r.dataSource.GetConnection(ctx)

	if err := db.AutoMigrate(&sobPO{}); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	return nil
}

func (r SobPostgresRepository) CreateSob(ctx context.Context, sob *sob.Sob) error {
	db := r.dataSource.GetConnection(ctx)

	po := sobBOToPO(*sob)
	return db.Create(&po).Error
}

func (r SobPostgresRepository) UpdateSob(ctx context.Context, sobId uuid.UUID, updateFn func(s *sob.Sob) (*sob.Sob, error)) error {
	db := r.dataSource.GetConnection(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		po := sobPO{Id: sobId}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po).Error; err != nil {
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

		po = sobBOToPO(*updatedBO)
		return tx.Save(&po).Error
	})
}

// Queries

func (r SobPostgresRepository) SearchSobs(
	ctx context.Context,
	pageRequest data.PageRequest,
) (data.Page[query.Sob], error) {
	return data.SearchEntities(ctx, pageRequest, sobPO{}, sobPOToDTO, r.dataSource.GetConnection(ctx))
}

func (r SobPostgresRepository) SobById(ctx context.Context, sobId uuid.UUID) (query.Sob, error) {
	db := r.dataSource.GetConnection(ctx)

	dbSob := sobPO{Id: sobId}
	if err := db.First(&dbSob).Error; err != nil {
		return query.Sob{}, fmt.Errorf("failed to read sob %s: %w", sobId, err)
	}

	return sobPOToDTO(dbSob), nil
}
