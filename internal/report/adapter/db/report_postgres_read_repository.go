package db

import (
	"context"
	"errors"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportPostgresReadRepository struct {
	dataSource datasource.DataSource
}

func NewReportPostgresReadRepository(dataSource datasource.DataSource) *ReportPostgresReadRepository {
	if dataSource == nil {
		panic("nil data source")
	}
	return &ReportPostgresReadRepository{dataSource: dataSource}
}

func (r ReportPostgresReadRepository) SearchReport(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Report], error) {
	addSobFilter(sobId, pageRequest)
	return data.SearchEntities(ctx, pageRequest, reportPO{}, reportPOToDTO, r.dataSource.GetConnection(ctx).
		Preload("Columns").
		Preload("Rows.Expression").
		Preload("Period"))
}

func (r ReportPostgresReadRepository) ReportById(ctx context.Context, reportId uuid.UUID) (query.Report, error) {
	var po reportPO
	err := r.dataSource.GetConnection(ctx).
		Preload("Columns").
		Preload("Rows.Expression").
		Preload("Period").
		Where("reports.id = ?", reportId).
		First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return query.Report{}, nil
	}
	if err != nil {
		return query.Report{}, err
	}
	return reportPOToDTO(po), nil
}

func (r ReportPostgresReadRepository) ReportInstanceByClassAndPeriod(ctx context.Context, sobId uuid.UUID, reportClass string, fiscalYear int, periodNumber int) (query.Report, error) {
	var po reportPO
	err := r.dataSource.GetConnection(ctx).
		Preload("Columns").
		Preload("Rows.Expression").
		Joins("Period").
		Where("reports.sob_id = ? AND reports.class = ? AND reports.template = false", sobId, reportClass).
		Where(`"Period".fiscal_year = ? AND "Period".period_number = ?`, fiscalYear, periodNumber).
		First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return query.Report{}, nil
	}
	if err != nil {
		return query.Report{}, err
	}
	return reportPOToDTO(po), nil
}

func (r ReportPostgresReadRepository) TemplateBySobIdAndClass(ctx context.Context, sobId uuid.UUID, reportClass string) (query.Report, error) {
	var po reportPO
	err := r.dataSource.GetConnection(ctx).
		Preload("Columns").
		Preload("Rows.Expression").
		Where("reports.sob_id = ? AND reports.class = ? AND reports.template = true", sobId, reportClass).
		First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return query.Report{}, nil
	}
	if err != nil {
		return query.Report{}, err
	}
	return reportPOToDTO(po), nil
}

func addSobFilter(sobId uuid.UUID, pageRequest data.PageRequest) {
	if sobId != uuid.Nil {
		sobIdFilter, _ := filterable.NewFilter("sobId", filterable.OptEq, sobId.String())
		pageRequest.AddAndFilterable(filterable.NewFilterableAtom(sobIdFilter))
	}
}
