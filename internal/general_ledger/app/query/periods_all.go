package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type AllPeriodsHandler struct {
	readModel GeneralLedgerReadModel
}

func NewAllPeriodsHandler(readModel GeneralLedgerReadModel) AllPeriodsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return AllPeriodsHandler{readModel: readModel}
}

func (h AllPeriodsHandler) Handle(ctx context.Context, sobId uuid.UUID) ([]Period, error) {
	sort1, err := sortable.NewSort("fiscalYear", "asc")
	if err != nil {
		panic(fmt.Errorf("failed to build sort 'fiscalYear': %w", err))
	}
	sort2, err := sortable.NewSort("periodNumber", "asc")
	if err != nil {
		panic(fmt.Errorf("failed to build sort 'periodNumber': %w", err))
	}

	pageRequest := data.NewPageRequest(
		pageable.Unpaged(),
		sortable.New(sort1, sort2),
		filterable.Unfiltered(),
	)
	periods, err := h.readModel.SearchPeriods(ctx, sobId, pageRequest)
	if err != nil {
		return nil, fmt.Errorf("error getting periods: %w", err)
	}
	return periods.Content(), nil
}
