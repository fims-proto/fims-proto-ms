package db

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/journal/domain/journal_entry"

	"github/fims-proto/fims-proto-ms/internal/common/data"

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

func (r JournalEntryPostgresRepository) JournalEntryById(ctx context.Context, entryId uuid.UUID) (query.JournalEntry, error) {
	db := readDBFromCtx(ctx)

	po := journalEntryPO{}
	if err := db.Preload("LineItems").First(&po, "entry_id = ?", entryId).Error; err != nil {
		return query.JournalEntry{}, errors.Wrap(err, "failed find journal entry by id")
	}

	return journalEntryPOToDTO(po), nil
}

func (r JournalEntryPostgresRepository) PagingJournalEntries(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[query.JournalEntry], error) {
	db := readDBFromCtx(ctx)

	var journalEntryPOs []journalEntryPO

	db.Scopes(data.Filtering(pageable)).Where("sob_id = ?", sobId)

	var count int64
	if err := db.Model(&journalEntryPO{}).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "count journal entries failed")
	}

	if err := db.Scopes(data.Paging(pageable)).Preload("LineItems").Find(&journalEntryPOs).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find journal entries by sobId %s", sobId)
	}

	var journalEntryDTOs []query.JournalEntry
	for _, po := range journalEntryPOs {
		journalEntryDTOs = append(journalEntryDTOs, journalEntryPOToDTO(po))
	}

	return data.NewPage(journalEntryDTOs, pageable, int(count))
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
