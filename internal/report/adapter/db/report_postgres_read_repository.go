package db

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"
)

type ReportPostgresReadRepository struct {
	dataSource datasource.DataSource
}

func NewReportPostgresReadRepository(dataSource datasource.DataSource) *ReportPostgresReadRepository {
	return &ReportPostgresReadRepository{
		dataSource: dataSource,
	}
}

func (r ReportPostgresReadRepository) SearchReport(
	ctx context.Context,
	sobId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[query.Report], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, reportPO{}, reportPOToDTO, r.dataSource.GetConnection(ctx).
		Preload("Sections.Items.Formulas.Account").
		Joins("Period"))
}

func addSobFilter(sobId uuid.UUID, pageRequest data.PageRequest) {
	if sobId != uuid.Nil {
		sobIdFilter, _ := filterable.NewFilter("sobId", filterable.OptEq, sobId.String())
		pageRequest.AddAndFilterable(filterable.NewFilterableAtom(sobIdFilter))
	}
}
