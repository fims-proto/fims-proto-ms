package db

import (
	"fmt"
	"strings"
	"time"

	"github/fims-proto/fims-proto-ms/internal/dimension/app/query"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain/category"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain/option"

	"github.com/google/uuid"
)

type dimensionCategoryPO struct {
	Id    uuid.UUID `gorm:"type:uuid;primaryKey"`
	SobId uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_DimCategories_SobId_Name"`
	Name  string    `gorm:"uniqueIndex:UQ_DimCategories_SobId_Name"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

type dimensionOptionPO struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey"`
	CategoryId uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_DimOptions_CategoryId_Name"`
	Name       string    `gorm:"uniqueIndex:UQ_DimOptions_CategoryId_Name"`

	Category dimensionCategoryPO `gorm:"foreignKey:CategoryId"`

	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// schemas — ResolveAssociation is required by the data.SearchEntities infrastructure

func (p dimensionCategoryPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return "dimension_categories", nil
	}

	return "", fmt.Errorf("dimensionCategoryPO doesn't have association named %s", entity)
}

func (p dimensionOptionPO) ResolveAssociation(entity string) (string, error) {
	if entity == "" {
		return "dimension_options", nil
	}

	if strings.EqualFold(entity, "category") {
		return "Category", nil
	}

	return "", fmt.Errorf("dimensionOptionPO doesn't have association named %s", entity)
}

// mappers: domain BO ↔ PO

func categoryBOToPO(bo *category.DimensionCategory) dimensionCategoryPO {
	return dimensionCategoryPO{
		Id:    bo.Id(),
		SobId: bo.SobId(),
		Name:  bo.Name(),
	}
}

func categoryPOToBO(po dimensionCategoryPO) (*category.DimensionCategory, error) {
	return category.New(po.Id, po.SobId, po.Name)
}

func optionBOToPO(bo *option.DimensionOption) dimensionOptionPO {
	return dimensionOptionPO{
		Id:         bo.Id(),
		CategoryId: bo.CategoryId(),
		Name:       bo.Name(),
	}
}

func optionPOToBO(po dimensionOptionPO) (*option.DimensionOption, error) {
	return option.New(po.Id, po.CategoryId, po.Name)
}

// mappers: PO → query DTO

func categoryPOToDTO(po dimensionCategoryPO) query.DimensionCategory {
	return query.DimensionCategory{
		Id:        po.Id,
		SobId:     po.SobId,
		Name:      po.Name,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}

func optionPOToDTO(po dimensionOptionPO) query.DimensionOption {
	return query.DimensionOption{
		Id:         po.Id,
		CategoryId: po.CategoryId,
		Name:       po.Name,
		Category:   categoryPOToDTO(po.Category),
		CreatedAt:  po.CreatedAt,
		UpdatedAt:  po.UpdatedAt,
	}
}
