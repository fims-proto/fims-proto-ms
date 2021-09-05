package db

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/counter/app/query"
	"github/fims-proto/fims-proto-ms/internal/counter/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CounterPostgresRepository struct{}

func NewCounterPostgresRepository() *CounterPostgresRepository {
	return &CounterPostgresRepository{}
}

func (r CounterPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	if err := db.AutoMigrate(&counter{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r CounterPostgresRepository) CreateCounter(ctx context.Context, c *domain.Counter) error {
	db := readDBFromCtx(ctx)

	dbCounter := marshall(c)

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(dbCounter).Error
	}); err != nil {
		return errors.Wrap(err, "create counter failed")
	}

	return nil
}

func (r CounterPostgresRepository) DeleteCounter(ctx context.Context, counterId uuid.UUID) error {
	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&counter{}, counterId).Error
	}); err != nil {
		return errors.Wrap(err, "delete counter failed")
	}

	return nil
}

func (r CounterPostgresRepository) UpdateCounter(
	ctx context.Context,
	counterId uuid.UUID,
	updateFn func(c *domain.Counter) (*domain.Counter, error),
) error {
	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		dbCounter := &counter{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(dbCounter, "id = ?", counterId).Error; err != nil {
			return err
		}

		domainCounter, err := unmarshallToDomain(dbCounter)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshall counter")
		}

		updatedDomainCounter, err := updateFn(domainCounter)
		if err != nil {
			return errors.Wrap(err, "failed to update counter")
		}

		dbCounter = marshall(updatedDomainCounter)
		if err := tx.Save(dbCounter).Error; err != nil {
			return errors.Wrap(err, "failed to save counter")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "update counter failed")
	}

	return nil
}

func (r CounterPostgresRepository) UpdateAndRead(
	ctx context.Context,
	counterId uuid.UUID,
	updateAndReadFn func(c *domain.Counter) (*domain.Counter, interface{}, error),
) (interface{}, error) {
	db := readDBFromCtx(ctx)

	var readValue interface{}

	if err := db.Transaction(func(tx *gorm.DB) error {
		dbCounter := &counter{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(dbCounter, "id = ?", counterId).Error; err != nil {
			return err
		}

		domainCounter, err := unmarshallToDomain(dbCounter)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshall counter")
		}

		var updatedDomainCounter *domain.Counter
		updatedDomainCounter, readValue, err = updateAndReadFn(domainCounter)
		if err != nil {
			return errors.Wrap(err, "failed to update counter in transaction")
		}

		dbCounter = marshall(updatedDomainCounter)
		if err := tx.Save(dbCounter).Error; err != nil {
			return errors.Wrap(err, "failed to save counter")
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "update counter failed")
	}

	return readValue, nil
}

func (r CounterPostgresRepository) ReadByBusinessObject(ctx context.Context, sep string, businessObjects []string) (query.Counter, error) {
	db := readDBFromCtx(ctx)

	m, err := domain.NewMatcher(sep, businessObjects...)
	if err != nil {
		return query.Counter{}, errors.Errorf("cannot create counter matcher: %s", businessObjects)
	}

	var dbCounter counter
	if err := db.Where("business_object = ?", m.String()).First(&dbCounter).Error; err != nil {
		return query.Counter{}, errors.Wrapf(err, "failed to read counter by business obejct %s", m.String())
	}
	return unmarshallToQuery(&dbCounter), nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
