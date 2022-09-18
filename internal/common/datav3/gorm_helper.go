package datav3

import (
	"context"

	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/schema"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/sortable"
	"gorm.io/gorm"
)

func SearchEntities[PO schema.Schema, DTO any](
	ctx context.Context,
	r PageRequest,
	po PO,
	convert func(po PO) (DTO, error),
	db *gorm.DB,
) (Page[DTO], error) {
	var persistentObjects []PO
	tx := db.Scopes(filterable.Filtering(r, po)).WithContext(ctx) // new session

	var count int64
	if err := tx.Model(&po).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count entities")
	}

	if err := tx.
		Scopes(pageable.Paging(r)).
		Scopes(sortable.Sorting(r, po)).
		Find(&persistentObjects).
		Error; err != nil {
		return nil, errors.Wrapf(err, "failed to search entities")
	}

	var dataTransferObjects []DTO
	for _, persistentObject := range persistentObjects {
		dto, err := convert(persistentObject)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map entity to DTO")
		}
		dataTransferObjects = append(dataTransferObjects, dto)
	}

	return NewPage(dataTransferObjects, r, int(count))
}
