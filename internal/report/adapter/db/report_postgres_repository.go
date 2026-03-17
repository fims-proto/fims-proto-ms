package db

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
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
	// First, delete all existing associations to ensure orphaned records are removed
	// This is necessary because GORM's FullSaveAssociations only upserts, it doesn't delete removed associations

	// Step 1: Find all section IDs for this report
	var sectionIds []uuid.UUID
	if err = db.Model(&sectionPO{}).Where("report_id = ?", reportId).Pluck("id", &sectionIds).Error; err != nil {
		return fmt.Errorf("failed to find section ids: %w", err)
	}

	// Step 2: Find all item IDs for these sections
	var itemIds []uuid.UUID
	if len(sectionIds) > 0 {
		if err = db.Model(&itemPO{}).Where("section_id IN ?", sectionIds).Pluck("id", &itemIds).Error; err != nil {
			return fmt.Errorf("failed to find item ids: %w", err)
		}
	}

	// Step 3: Delete formulas associated with these items
	if len(itemIds) > 0 {
		if err = db.Where("item_id IN ?", itemIds).Delete(&formulaPO{}).Error; err != nil {
			return fmt.Errorf("failed to delete formulas: %w", err)
		}
	}

	// Step 4: Delete items associated with these sections
	if len(sectionIds) > 0 {
		if err = db.Where("section_id IN ?", sectionIds).Delete(&itemPO{}).Error; err != nil {
			return fmt.Errorf("failed to delete items: %w", err)
		}
	}

	// Step 5: Delete all sections for this report
	if err = db.Where("report_id = ?", reportId).Delete(&sectionPO{}).Error; err != nil {
		return fmt.Errorf("failed to delete sections: %w", err)
	}

	// Step 6: Save the updated report with all new associations
	updatedPO := reportBOToPO(updatedBO)
	return db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&updatedPO).Error
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
