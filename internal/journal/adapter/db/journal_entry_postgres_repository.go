package db

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/datav3"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/sortable"

	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/journal/app/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JournalEntryPostgresRepository struct{}

func NewJournalEntryPostgresRepository() *JournalEntryPostgresRepository {
	return &JournalEntryPostgresRepository{}
}

func (r JournalEntryPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	if err := db.AutoMigrate(&journalEntryPO{}, &lineItemPO{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r JournalEntryPostgresRepository) CreateJournalEntry(ctx context.Context, d *journal_entry.JournalEntry) error {
	db := readDBFromCtx(ctx)

	po := journalEntryBOToPO(*d)

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&po).Error
	})
}

func (r JournalEntryPostgresRepository) UpdateJournalEntry(ctx context.Context, entryId uuid.UUID, updateFn func(d *journal_entry.JournalEntry) (*journal_entry.JournalEntry, error)) error {
	db := readDBFromCtx(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		po := journalEntryPO{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("LineItems").First(&po, "entry_id = ?", entryId).Error; err != nil {
			return err
		}

		bo, err := journalEntryPOToBO(po)
		if err != nil {
			return errors.Wrap(err, "failed to map journal entry")
		}

		updatedBO, err := updateFn(bo)
		if err != nil {
			return errors.Wrap(err, "update journal entry in transaction failed")
		}

		po = journalEntryBOToPO(*updatedBO)

		// remove existing first
		if err = tx.Where("entry_id = ?", po.EntryId).Delete(&lineItemPO{}).Error; err != nil {
			return errors.Wrap(err, "delete journal entry items failed")
		}

		return tx.Save(&po).Error
	})
}

// queries

func (r JournalEntryPostgresRepository) SearchJournalEntries(ctx context.Context, sobId uuid.UUID, pageRequest datav3.PageRequest) (datav3.Page[query.JournalEntry], error) {
	if sobId != uuid.Nil {
		sobIdFilter, _ := filterable.NewFilter("sobId", "eq", sobId.String())
		pageRequest.AddFilter(sobIdFilter)
	}
	return datav3.SearchEntities(ctx, pageRequest, journalEntryPO{}, journalEntryPOToDTO, readDBFromCtx(ctx).Preload("LineItems"))
}

func (r JournalEntryPostgresRepository) JournalEntryById(ctx context.Context, entryId uuid.UUID) (query.JournalEntry, error) {
	entryIdFilter, _ := filterable.NewFilter("entryId", "eq", entryId)
	pageRequest := datav3.NewPageRequest(pageable.Unpaged(), sortable.Unsorted(), filterable.New(entryIdFilter))

	journalEntries, err := r.SearchJournalEntries(ctx, uuid.Nil, pageRequest)
	if err != nil {
		return query.JournalEntry{}, err
	}

	if journalEntries.NumberOfElements() != 1 {
		return query.JournalEntry{}, errors.Errorf("journal entry not found by id: %s", entryId)
	}

	return journalEntries.Content()[0], nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
