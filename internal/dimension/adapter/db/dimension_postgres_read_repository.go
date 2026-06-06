package db

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/dimension/app/query"

	"github.com/google/uuid"
)

type DimensionPostgresReadRepository struct {
	dataSource datasource.DataSource
}

func NewDimensionPostgresReadRepository(dataSource datasource.DataSource) *DimensionPostgresReadRepository {
	if dataSource == nil {
		panic("nil data source")
	}

	return &DimensionPostgresReadRepository{dataSource: dataSource}
}

func (r DimensionPostgresReadRepository) SearchCategories(
	ctx context.Context,
	sobId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.DimensionCategory], error) {
	addSobFilter(sobId, pageRequest)

	return data.SearchEntities(ctx, pageRequest, dimensionCategoryPO{}, categoryPOToDTO, r.dataSource.GetConnection(ctx))
}

func (r DimensionPostgresReadRepository) SearchOptions(
	ctx context.Context,
	categoryId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.DimensionOption], error) {
	addCategoryFilter(categoryId, pageRequest)

	return data.SearchEntities(ctx, pageRequest, dimensionOptionPO{}, optionPOToDTO, r.dataSource.GetConnection(ctx))
}

func (r DimensionPostgresReadRepository) CategoryById(
	ctx context.Context,
	categoryId uuid.UUID,
) (query.DimensionCategory, error) {
	db := r.dataSource.GetConnection(ctx)

	var po dimensionCategoryPO
	if err := db.First(&po, categoryId).Error; err != nil {
		if isNotFound(err) {
			return query.DimensionCategory{}, commonErrors.ErrRecordNotFound()
		}

		return query.DimensionCategory{}, err
	}

	return categoryPOToDTO(po), nil
}

func (r DimensionPostgresReadRepository) OptionsByIds(
	ctx context.Context,
	optionIds []uuid.UUID,
) ([]query.DimensionOption, error) {
	if len(optionIds) == 0 {
		return nil, nil
	}

	db := r.dataSource.GetConnection(ctx)

	var pos []dimensionOptionPO
	if err := db.Where("id IN ?", optionIds).Preload("Category").Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("failed to read dimension options: %w", err)
	}

	result := make([]query.DimensionOption, 0, len(pos))
	for _, po := range pos {
		result = append(result, optionPOToDTO(po))
	}

	return result, nil
}

func (r DimensionPostgresReadRepository) CategoriesByIds(
	ctx context.Context,
	categoryIds []uuid.UUID,
) ([]query.DimensionCategory, error) {
	if len(categoryIds) == 0 {
		return nil, nil
	}

	db := r.dataSource.GetConnection(ctx)

	var pos []dimensionCategoryPO
	if err := db.Where("id IN ?", categoryIds).Find(&pos).Error; err != nil {
		return nil, fmt.Errorf("failed to read dimension categories: %w", err)
	}

	result := make([]query.DimensionCategory, 0, len(pos))
	for _, po := range pos {
		result = append(result, categoryPOToDTO(po))
	}

	return result, nil
}

// helpers

func addSobFilter(sobId uuid.UUID, pageRequest data.PageRequest) {
	if sobId != uuid.Nil {
		sobIdFilter, _ := filterable.NewFilter("sobId", filterable.OptEq, sobId.String())
		pageRequest.AddAndFilterable(filterable.NewFilterableAtom(sobIdFilter))
	}
}

func addCategoryFilter(categoryId uuid.UUID, pageRequest data.PageRequest) {
	if categoryId != uuid.Nil {
		filter, _ := filterable.NewFilter("categoryId", filterable.OptEq, categoryId.String())
		pageRequest.AddAndFilterable(filterable.NewFilterableAtom(filter))
	}
}
