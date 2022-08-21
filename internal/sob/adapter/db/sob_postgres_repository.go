package db

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/sob/domain/sob"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SobPostgresRepository struct{}

func NewSobPostgresRepository() *SobPostgresRepository {
	return &SobPostgresRepository{}
}

func (r SobPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	if err := db.AutoMigrate(&sobPO{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r SobPostgresRepository) CreateSob(ctx context.Context, sob *sob.Sob) error {
	db := readDBFromCtx(ctx)

	po, err := sobBOToPO(*sob)
	if err != nil {
		return errors.Wrap(err, "failed to sobBOToPO sob")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&po).Error
	})
}

func (r SobPostgresRepository) UpdateSob(ctx context.Context, sobId uuid.UUID, updateFn func(s *sob.Sob) (*sob.Sob, error)) error {
	db := readDBFromCtx(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		po := sobPO{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po, "id = ?", sobId).Error; err != nil {
			return errors.Wrap(err, "failed to find sob")
		}

		bo, err := sobPOToBO(po)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal sob")
		}

		updatedBO, err := updateFn(bo)
		if err != nil {
			return errors.Wrap(err, "failed to update sob in transaction")
		}

		po, err = sobBOToPO(*updatedBO)
		if err != nil {
			return errors.Wrap(err, "failed to sobBOToPO sob")
		}

		return tx.Save(&po).Error
	})
}

// Queries

func (r SobPostgresRepository) PagingSobs(ctx context.Context, pageable data.Pageable) (data.Page[query.Sob], error) {
	db := readDBFromCtx(ctx)

	var sobPOs []sobPO

	db = db.Scopes(data.Filtering(pageable))

	var count int64
	if err := db.Model(&sobPO{}).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count sobs")
	}

	if err := db.Scopes(data.Paging(pageable)).Find(&sobPOs).Error; err != nil {
		return nil, errors.Wrap(err, "failed to find sobs")
	}

	var querySobs []query.Sob
	for _, dbSob := range sobPOs {
		querySob, err := sobPOToDTO(dbSob)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal sob")
		}
		querySobs = append(querySobs, querySob)
	}

	return data.NewPage(querySobs, pageable, int(count))
}

func (r SobPostgresRepository) SobById(ctx context.Context, sobId uuid.UUID) (query.Sob, error) {
	db := readDBFromCtx(ctx)

	dbSob := sobPO{}
	if err := db.Where("id = ?", sobId).First(&dbSob).Error; err != nil {
		return query.Sob{}, errors.Wrapf(err, "failed to read sob %s", sobId)
	}

	querySob, err := sobPOToDTO(dbSob)
	if err != nil {
		return query.Sob{}, errors.Wrap(err, "failed to unmarshal sob")
	}

	return querySob, nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
