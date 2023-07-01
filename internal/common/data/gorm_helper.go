package data

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/schema"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"
	"gorm.io/gorm"
)

func SearchEntities[PO schema.Schema, DTO any](
	ctx context.Context,
	r PageRequest,
	po PO,
	convert func(po PO) DTO,
	db *gorm.DB,
) (Page[DTO], error) {
	var persistentObjects []PO
	tx := db.Scopes(filterable.Filtering(r, po)).Session(&gorm.Session{}) // new session

	var count int64
	if err := tx.Model(&po).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("failed to count entities: %w", err)
	}

	if err := tx.
		Scopes(pageable.Paging(r)).
		Scopes(sortable.Sorting(r, po)).
		Find(&persistentObjects).
		Error; err != nil {
		return nil, fmt.Errorf("failed to search entities: %w", err)
	}

	var dataTransferObjects []DTO
	for _, persistentObject := range persistentObjects {
		dto := convert(persistentObject)
		dataTransferObjects = append(dataTransferObjects, dto)
	}

	return NewPage(dataTransferObjects, r, int(count))
}
