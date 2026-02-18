package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

type AuxiliaryCategoryByKeyHandler struct {
	readModel GeneralLedgerReadModel
}

func NewAuxiliaryCategoryByKeyHandler(readModel GeneralLedgerReadModel) AuxiliaryCategoryByKeyHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return AuxiliaryCategoryByKeyHandler{
		readModel: readModel,
	}
}

func (h AuxiliaryCategoryByKeyHandler) Handle(ctx context.Context, sobId uuid.UUID, categoryKey string) (AuxiliaryCategory, error) {
	keyFilter, err := filterable.NewFilter("key", filterable.OptEq, categoryKey)
	if err != nil {
		panic("failed to create filter 'key'")
	}

	categories, err := h.readModel.SearchAuxiliaryCategories(ctx, sobId, data.NewPageRequest(pageable.Unpaged(), sortable.Unsorted(), filterable.NewFilterableAtom(keyFilter)))
	if err != nil {
		return AuxiliaryCategory{}, fmt.Errorf("failed to search auxiliary categories: %w", err)
	}
	if categories.NumberOfElements() != 1 {
		return AuxiliaryCategory{}, commonErrors.ErrRecordNotFound()
	}

	return categories.Content()[0], nil
}
