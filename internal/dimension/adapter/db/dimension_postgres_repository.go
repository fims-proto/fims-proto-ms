package db

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain/category"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain/option"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DimensionPostgresRepository struct {
	dataSource datasource.DataSource
}

func NewDimensionPostgresRepository(dataSource datasource.DataSource) *DimensionPostgresRepository {
	if dataSource == nil {
		panic("nil data source")
	}

	return &DimensionPostgresRepository{dataSource: dataSource}
}

func (r DimensionPostgresRepository) Migrate(ctx context.Context) error {
	db := r.dataSource.GetConnection(ctx)

	return db.AutoMigrate(
		&dimensionCategoryPO{},
		&dimensionOptionPO{},
	)
}

func (r DimensionPostgresRepository) EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error {
	return r.dataSource.EnableTransaction(ctx, txFn)
}

// Category operations

func (r DimensionPostgresRepository) CreateCategory(ctx context.Context, c *category.DimensionCategory) error {
	db := r.dataSource.GetConnection(ctx)

	po := categoryBOToPO(c)

	return db.Create(&po).Error
}

func (r DimensionPostgresRepository) UpdateCategory(
	ctx context.Context,
	categoryId uuid.UUID,
	updateFn func(c *category.DimensionCategory) (*category.DimensionCategory, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	po := dimensionCategoryPO{Id: categoryId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&po).Error; err != nil {
		return err
	}

	bo, err := categoryPOToBO(po)
	if err != nil {
		return fmt.Errorf("failed to load dimension category: %w", err)
	}

	updatedBO, err := updateFn(bo)
	if err != nil {
		return fmt.Errorf("failed to update dimension category: %w", err)
	}

	po = categoryBOToPO(updatedBO)

	return db.Save(&po).Error
}

func (r DimensionPostgresRepository) DeleteCategory(ctx context.Context, categoryId uuid.UUID) error {
	db := r.dataSource.GetConnection(ctx)

	return db.Delete(&dimensionCategoryPO{Id: categoryId}).Error
}

func (r DimensionPostgresRepository) ReadCategoryById(ctx context.Context, categoryId uuid.UUID) (*category.DimensionCategory, error) {
	db := r.dataSource.GetConnection(ctx)

	var po dimensionCategoryPO
	if err := db.First(&po, categoryId).Error; err != nil {
		if isNotFound(err) {
			return nil, commonErrors.ErrRecordNotFound()
		}

		return nil, err
	}

	return categoryPOToBO(po)
}

func (r DimensionPostgresRepository) ExistsCategoryUsedByJournalLine(ctx context.Context, categoryId uuid.UUID) (bool, error) {
	db := r.dataSource.GetConnection(ctx)

	var count int64
	err := db.Model(&dimensionOptionPO{}).
		Joins("JOIN journal_line_dimension_options jldo ON jldo.dimension_option_id = dimension_options.id").
		Where("dimension_options.category_id = ?", categoryId).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check category usage: %w", err)
	}

	return count > 0, nil
}

// Option operations

func (r DimensionPostgresRepository) CreateOption(ctx context.Context, o *option.DimensionOption) error {
	db := r.dataSource.GetConnection(ctx)

	po := optionBOToPO(o)

	return db.Create(&po).Error
}

func (r DimensionPostgresRepository) UpdateOption(
	ctx context.Context,
	optionId uuid.UUID,
	updateFn func(o *option.DimensionOption) (*option.DimensionOption, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	po := dimensionOptionPO{Id: optionId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&po).Error; err != nil {
		return err
	}

	bo, err := optionPOToBO(po)
	if err != nil {
		return fmt.Errorf("failed to load dimension option: %w", err)
	}

	updatedBO, err := updateFn(bo)
	if err != nil {
		return fmt.Errorf("failed to update dimension option: %w", err)
	}

	po = optionBOToPO(updatedBO)

	return db.Save(&po).Error
}

func (r DimensionPostgresRepository) DeleteOption(ctx context.Context, optionId uuid.UUID) error {
	db := r.dataSource.GetConnection(ctx)

	return db.Delete(&dimensionOptionPO{Id: optionId}).Error
}

func (r DimensionPostgresRepository) DeleteOptionsByCategoryId(ctx context.Context, categoryId uuid.UUID) error {
	db := r.dataSource.GetConnection(ctx)

	return db.Where("category_id = ?", categoryId).Delete(&dimensionOptionPO{}).Error
}

func (r DimensionPostgresRepository) ReadOptionsByIds(ctx context.Context, optionIds []uuid.UUID) ([]*option.DimensionOption, error) {
	if len(optionIds) == 0 {
		return nil, nil
	}

	db := r.dataSource.GetConnection(ctx)

	var pos []dimensionOptionPO
	if err := db.Where("id IN ?", optionIds).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("failed to read dimension options: %w", err)
	}

	result := make([]*option.DimensionOption, 0, len(pos))

	for _, po := range pos {
		bo, err := optionPOToBO(po)
		if err != nil {
			return nil, err
		}

		result = append(result, bo)
	}

	return result, nil
}

func (r DimensionPostgresRepository) ExistsOptionUsedByJournalLine(ctx context.Context, optionId uuid.UUID) (bool, error) {
	db := r.dataSource.GetConnection(ctx)

	var count int64
	err := db.Table("journal_line_dimension_options").
		Where("dimension_option_id = ?", optionId).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check option usage: %w", err)
	}

	return count > 0, nil
}

// helpers

func isNotFound(err error) bool {
	return err != nil && err.Error() == gorm.ErrRecordNotFound.Error()
}
