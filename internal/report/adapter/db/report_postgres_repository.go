package db

import (
	"context"
	"errors"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReportPostgresRepository struct {
	dataSource datasource.DataSource
}

func NewReportPostgresRepository(dataSource datasource.DataSource) *ReportPostgresRepository {
	if dataSource == nil {
		panic("nil data source")
	}
	return &ReportPostgresRepository{dataSource: dataSource}
}

func (r ReportPostgresRepository) Migrate(ctx context.Context) error {
	db := r.dataSource.GetConnection(ctx)

	return db.AutoMigrate(
		&reportPO{},
		&reportColumnPO{},
		&reportRowPO{},
		&reportExpressionPO{},
	)
}

func (r ReportPostgresRepository) EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error {
	return r.dataSource.EnableTransaction(ctx, txFn)
}

func (r ReportPostgresRepository) CreateReports(ctx context.Context, reports []*report.Report) error {
	db := r.dataSource.GetConnection(ctx)
	pos := converter.BOsToPOs(reports, reportBOToPO)
	return commonErrors.TranslateDBError(db.Create(pos).Error)
}

func (r ReportPostgresRepository) UpdateReport(
	ctx context.Context,
	reportId uuid.UUID,
	updateFn func(r *report.Report) (*report.Report, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	po := reportPO{Id: reportId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Columns").
		Preload("Rows.Expression").
		Preload("Period").
		First(&po).Error; err != nil {
		return err
	}

	bo, err := reportPOToBO(&po)
	if err != nil {
		return err
	}
	updatedBO, err := updateFn(bo)
	if err != nil {
		return err
	}

	if err = r.deleteReportAssociations(ctx, reportId); err != nil {
		return err
	}

	updatedPO := reportBOToPO(updatedBO)
	return commonErrors.TranslateDBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&updatedPO).Error)
}

func (r ReportPostgresRepository) deleteReportAssociations(ctx context.Context, reportId uuid.UUID) error {
	db := r.dataSource.GetConnection(ctx)

	var rowIds []uuid.UUID
	if err := db.Model(&reportRowPO{}).Where("report_id = ?", reportId).Pluck("id", &rowIds).Error; err != nil {
		return fmt.Errorf("failed to find report row ids: %w", err)
	}

	if len(rowIds) > 0 {
		if err := db.Where("row_id IN ?", rowIds).Delete(&reportExpressionPO{}).Error; err != nil {
			return fmt.Errorf("failed to delete report expressions: %w", err)
		}
		if err := db.Where("id IN ?", rowIds).Delete(&reportRowPO{}).Error; err != nil {
			return fmt.Errorf("failed to delete report rows: %w", err)
		}
	}

	if err := db.Where("report_id = ?", reportId).Delete(&reportColumnPO{}).Error; err != nil {
		return fmt.Errorf("failed to delete report columns: %w", err)
	}

	return nil
}

func (r ReportPostgresRepository) ReadReportById(ctx context.Context, reportId uuid.UUID) (*report.Report, error) {
	db := r.dataSource.GetConnection(ctx)
	po := reportPO{Id: reportId}
	if err := db.Preload("Columns").
		Preload("Rows.Expression").
		Preload("Period").
		First(&po).Error; err != nil {
		return nil, err
	}
	return reportPOToBO(&po)
}

func (r ReportPostgresRepository) ReadTemplatesBySobId(ctx context.Context, sobId uuid.UUID) ([]*report.Report, error) {
	db := r.dataSource.GetConnection(ctx)
	var pos []*reportPO
	if err := db.Preload("Columns").
		Preload("Rows.Expression").
		Where("sob_id = ? AND template = true", sobId).
		Find(&pos).Error; err != nil {
		return nil, err
	}
	return converter.POsToBOs(pos, reportPOToBO)
}

func (r ReportPostgresRepository) ReadTemplateBySobIdAndClass(ctx context.Context, sobId uuid.UUID, reportClass string) (*report.Report, error) {
	db := r.dataSource.GetConnection(ctx)
	var po reportPO
	err := db.Preload("Columns").
		Preload("Rows.Expression").
		Where("sob_id = ? AND template = true AND class = ?", sobId, reportClass).
		First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return reportPOToBO(&po)
}

func (r ReportPostgresRepository) ReadInstanceBySobClassAndPeriod(ctx context.Context, sobId uuid.UUID, reportClass string, periodId uuid.UUID) (*report.Report, error) {
	db := r.dataSource.GetConnection(ctx)
	var po reportPO
	err := db.Preload("Columns").
		Preload("Rows.Expression").
		Preload("Period").
		Where("sob_id = ? AND class = ? AND period_id = ? AND template = false", sobId, reportClass, periodId).
		First(&po).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return reportPOToBO(&po)
}
