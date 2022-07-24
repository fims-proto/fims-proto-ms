package db

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"

	"github.com/google/uuid"

	"gorm.io/gorm/clause"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type LedgerPostgresRepository struct{}

func NewLedgerPostgresRepository() *LedgerPostgresRepository {
	return &LedgerPostgresRepository{}
}

func (r LedgerPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	if err := db.AutoMigrate(&period{}, &ledger{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r LedgerPostgresRepository) CreatePeriod(ctx context.Context, period *domain.Period) error {
	db := readDBFromCtx(ctx)

	dbPeriod := marshalPeriod(*period)

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(dbPeriod).Error
	}); err != nil {
		return errors.Wrap(err, "failed to create period")
	}
	return nil
}

func (r LedgerPostgresRepository) CreateLedgers(ctx context.Context, ledgers []*domain.Ledger) error {
	db := readDBFromCtx(ctx)

	var dbLedgers []ledger
	for _, domainLedger := range ledgers {
		dbLedgers = append(dbLedgers, marshalLedger(*domainLedger))
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.CreateInBatches(dbLedgers, 100).Error
	}); err != nil {
		return errors.Wrap(err, "failed to create ledgers")
	}
	return nil
}

func (r LedgerPostgresRepository) UpdateLedgersByPeriodAndAccounts(ctx context.Context, periodId uuid.UUID, accountIds []uuid.UUID, updateFn func(ledgers []*domain.Ledger) ([]*domain.Ledger, error)) error {
	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		var dbLedgers []ledger
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&dbLedgers, "period_id = ? AND account_id IN ?", periodId, accountIds).Error; err != nil {
			return err
		}
		return r.updateLedgers(tx, dbLedgers, updateFn)
	}); err != nil {
		return errors.Wrap(err, "update ledger failed")
	}
	return nil
}

func (r LedgerPostgresRepository) updateLedgers(tx *gorm.DB, dbLedgers []ledger, updateFn func(ledgers []*domain.Ledger) ([]*domain.Ledger, error)) error {
	var ledgers []*domain.Ledger
	for _, dbLedger := range dbLedgers {
		domainLedger, err := unmarshalLedgerToDomain(dbLedger)
		if err != nil {
			return errors.Wrap(err, "unmarshal ledger failed")
		}
		ledgers = append(ledgers, domainLedger)
	}

	updatedLedgers, err := updateFn(ledgers)
	if err != nil {
		return errors.Wrap(err, "update ledger in transaction failed")
	}

	dbLedgers = nil // empty slice
	for _, updatedLedger := range updatedLedgers {
		dbLedgers = append(dbLedgers, marshalLedger(*updatedLedger))
	}
	if err := tx.Save(&dbLedgers).Error; err != nil {
		return errors.Wrap(err, "save ledger failed")
	}
	return nil
}

func (r LedgerPostgresRepository) ReadLedgersByPeriod(ctx context.Context, periodId uuid.UUID, pageable data.Pageable) (data.Page[query.Ledger], error) {
	db := readDBFromCtx(ctx)

	var dbLedgers []ledger

	db = data.AddFilter(pageable, db).Where("period_id = ?", periodId)

	var count int64
	if err := db.Model(&ledger{}).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "count ledgers failed")
	}

	db = data.AddPaging(pageable, db)

	if err := db.Find(&dbLedgers).Error; err != nil {
		return nil, errors.Wrapf(err, "find ledgers by period %s failed", periodId)
	}

	queryLedgers := make([]query.Ledger, len(dbLedgers))
	for i, dbLedger := range dbLedgers {
		queryLedgers[i] = unmarshalLedgerToQuery(dbLedger)
	}
	return data.NewPage(queryLedgers, pageable, int(count))
}

func (r LedgerPostgresRepository) ReadPeriodById(ctx context.Context, id uuid.UUID) (query.Period, error) {
	db := readDBFromCtx(ctx)

	dbPeriod := period{}
	if err := db.First(&dbPeriod, "id = ?", id).Error; err != nil {
		return query.Period{}, errors.Wrap(err, "find period by id failed")
	}

	return unmarshalPeriodToQuery(dbPeriod), nil
}

func (r LedgerPostgresRepository) ReadPeriodsByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]query.Period, error) {
	db := readDBFromCtx(ctx)

	var dbPeriods []period
	if err := db.Find(&dbPeriods, "id IN ?", ids).Error; err != nil {
		return nil, errors.Wrap(err, "find period by id failed")
	}

	periods := make(map[uuid.UUID]query.Period)
	for _, dbPeriod := range dbPeriods {
		periods[dbPeriod.Id] = unmarshalPeriodToQuery(dbPeriod)
	}

	return periods, nil
}

func (r LedgerPostgresRepository) ReadPeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (query.Period, error) {
	db := readDBFromCtx(ctx)

	var dbPeriods []period
	if err := db.Where("sob_id = ? AND opening_time <= ? AND (ending_time > ? OR ending_time = ?)", sobId, timePoint, timePoint, time.Time{}).Find(&dbPeriods).Error; err != nil {
		return query.Period{}, errors.Wrap(err, "find period by id failed")
	}

	if len(dbPeriods) != 1 {
		return query.Period{}, errors.Errorf("expected 1 but %d periods found", len(dbPeriods))
	}

	return unmarshalPeriodToQuery(dbPeriods[0]), nil
}

func (r LedgerPostgresRepository) ReadPeriods(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[query.Period], error) {
	db := readDBFromCtx(ctx)

	var dbPeriods []period

	db = data.AddFilter(pageable, db).Where("sob_id = ?", sobId)

	var count int64
	if err := db.Model(&period{}).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "count periods failed")
	}

	db = data.AddPaging(pageable, db)

	if err := db.Find(&dbPeriods).Error; err != nil {
		return nil, errors.Wrapf(err, "find periods by sob %s failed", sobId)
	}

	var queryPeriods []query.Period
	for _, dbPeriod := range dbPeriods {
		queryPeriods = append(queryPeriods, unmarshalPeriodToQuery(dbPeriod))
	}
	return data.NewPage(queryPeriods, pageable, int(count))
}

func (r LedgerPostgresRepository) ReadOpenPeriod(ctx context.Context, sobId uuid.UUID) (query.Period, error) {
	db := readDBFromCtx(ctx)

	var dbPeriods []period
	if err := db.Where("sob_id = ? AND is_closed = false", sobId).Find(&dbPeriods).Error; err != nil {
		return query.Period{}, errors.Wrapf(err, "find open period by sob %s failed", sobId)
	}

	if len(dbPeriods) != 1 {
		return query.Period{}, errors.Errorf("expects 1 open period, but find %d", len(dbPeriods))
	}

	return unmarshalPeriodToQuery(dbPeriods[0]), nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
