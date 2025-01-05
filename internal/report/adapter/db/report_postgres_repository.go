package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"gorm.io/gorm/clause"
)

type ReportPostgresRepository struct {
	dataSource datasource.DataSource
}

func NewReportPostgresRepository(dataSource datasource.DataSource) *ReportPostgresRepository {
	if dataSource == nil {
		panic("nil data source")
	}

	return &ReportPostgresRepository{
		dataSource: dataSource,
	}
}

func (r ReportPostgresRepository) Migrate(ctx context.Context) error {
	db := r.dataSource.GetConnection(ctx)

	return db.AutoMigrate(
		&reportPO{},
		&sectionPO{},
		&itemPO{},
		&formulaPO{},
	)
}

func (r ReportPostgresRepository) EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error {
	return r.dataSource.EnableTransaction(ctx, txFn)
}

func (r ReportPostgresRepository) CreateReports(ctx context.Context, reports []*report.Report) error {
	db := r.dataSource.GetConnection(ctx)

	pos := converter.BOsToPOs(reports, reportBOToPO)

	return db.Create(pos).Error
}

func (r ReportPostgresRepository) UpdateReport(
	ctx context.Context,
	reportId uuid.UUID,
	updateFn func(r *report.Report) (*report.Report, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	// read report
	po := reportPO{Id: reportId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Sections.Items.Formulas.Account").
		Preload("Period").
		First(&po).Error; err != nil {
		return err
	}

	bo, err := reportPOToBO(&po)
	if err != nil {
		return err
	}

	// delegate update
	updatedBO, err := updateFn(bo)
	if err != nil {
		return err
	}

	// save
	updatedPO := reportBOToPO(updatedBO)
	// TODO even the formula amounts in the updatedPO is correct, seems the db.Save cannot update the nested object.
	// TODO this could also happen to items and sections
	// TODO maybe need to updated recursively from bottom to top?
	return db.Save(&updatedPO).Error
}

func (r ReportPostgresRepository) ReadReportById(ctx context.Context, reportId uuid.UUID) (*report.Report, error) {
	db := r.dataSource.GetConnection(ctx)

	po := reportPO{Id: reportId}
	if err := db.Preload("Sections.Items.Formulas.Account").
		Joins("Period").
		First(&po).Error; err != nil {
		return nil, err
	}

	return reportPOToBO(&po)
}

func (r ReportPostgresRepository) UpdateItem(
	ctx context.Context,
	itemId uuid.UUID,
	updateFn func(i *report.Item) (*report.Item, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	// read item
	po := itemPO{Id: itemId}
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Formulas.Account").
		First(&po).Error; err != nil {
		return err
	}

	bo, err := itemPOToBO(&po)
	if err != nil {
		return err
	}

	// delegate update
	updatedBO, err := updateFn(bo)
	if err != nil {
		return err
	}

	// save
	// use the section id from the original po
	updatedPO := itemBOToPO(updatedBO, po.SectionId)
	// delete formulas first
	if err = db.Where("item_id = ?", updatedPO.Id).Delete(&formulaPO{}).Error; err != nil {
		return fmt.Errorf("failed to delete formulas: %w", err)
	}
	return db.Save(&updatedPO).Error
}
