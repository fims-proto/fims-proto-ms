package datav3

import (
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/sortable"
	"gorm.io/gorm"
)

func SearchEntities[PO any, DTO any](
	r PageRequest,
	singlePO *PO,
	convert func(po PO) (DTO, error),
	resolveEntity func(entity string) (string, error),
	db *gorm.DB,
) (Page[DTO], error) {
	var persistentObjects []PO
	db = db.Scopes(filterable.Filtering(r, resolveEntity))

	var count int64
	if err := db.Model(singlePO).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count entities")
	}

	if err := db.
		Scopes(pageable.Paging(r)).
		Scopes(sortable.Sorting(r, resolveEntity)).
		Find(&persistentObjects).
		Error; err != nil {
		return nil, errors.Wrapf(err, "failed to search entities")
	}

	var dataTransferObjects []DTO
	for _, po := range persistentObjects {
		dto, err := convert(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map entity to DTO")
		}
		dataTransferObjects = append(dataTransferObjects, dto)
	}

	return NewPage(dataTransferObjects, r, int(count))
}
