package db

import (
	"context"
	"time"

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

	if err := db.AutoMigrate(&accountingPeriod{}, &ledger{}, &ledgerLog{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r LedgerPostgresRepository) CreateAccountingPeriod(ctx context.Context, period *domain.AccountingPeriod) error {
	db := readDBFromCtx(ctx)

	dbPeriod := marshallAccountingPeriod(period)

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(dbPeriod).Error
	}); err != nil {
		return errors.Wrap(err, "failed to create accounting period")
	}
	return nil
}

func (r LedgerPostgresRepository) UpdateAccountingPeriod(ctx context.Context, id uuid.UUID, updateFn func(period *domain.AccountingPeriod) (*domain.AccountingPeriod, error)) error {
	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		dbPeriod := &accountingPeriod{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(dbPeriod, "id = ?", id).Error; err != nil {
			return err
		}

		period, err := unmarshallPeriodToDomain(dbPeriod)
		if err != nil {
			return errors.Wrap(err, "unmarshall accounting period failed")
		}

		updatedPeriod, err := updateFn(period)
		if err != nil {
			return errors.Wrap(err, "update accounting period in transaction failed")
		}

		dbPeriod = marshallAccountingPeriod(updatedPeriod)
		if err := tx.Save(dbPeriod).Error; err != nil {
			return errors.Wrap(err, "save accounting period failed")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "update accounting period failed")
	}
	return nil
}

func (r LedgerPostgresRepository) CreateLedgers(ctx context.Context, ledgers []*domain.Ledger) error {
	db := readDBFromCtx(ctx)

	var dbLedgers []*ledger
	for _, ledger := range ledgers {
		dbLedgers = append(dbLedgers, marshallLedger(ledger))
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.CreateInBatches(dbLedgers, 100).Error
	}); err != nil {
		return errors.Wrap(err, "failed to create ledgers")
	}
	return nil
}

func (r LedgerPostgresRepository) UpdateLedgers(ctx context.Context, ids []uuid.UUID, updateFn func(ledgers []*domain.Ledger) ([]*domain.Ledger, error)) error {
	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		var dbLedgers []ledger
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Find(&dbLedgers, "id IN ?", ids).Error; err != nil {
			return err
		}

		var ledgers []*domain.Ledger
		for _, dbLedger := range dbLedgers {
			ledger, err := unmarshallLedgerToDomain(&dbLedger)
			if err != nil {
				return errors.Wrap(err, "unmarshall ledger failed")
			}
			ledgers = append(ledgers, ledger)
		}

		updatedLedgers, err := updateFn(ledgers)
		if err != nil {
			return errors.Wrap(err, "update ledger in transaction failed")
		}

		dbLedgers = nil // empty slice
		for _, updatedLedger := range updatedLedgers {
			dbLedgers = append(dbLedgers, *marshallLedger(updatedLedger))
		}
		if err := tx.Save(&dbLedgers).Error; err != nil {
			return errors.Wrap(err, "save ledger failed")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "update ledger failed")
	}
	return nil
}

func (r LedgerPostgresRepository) CreateLedgerLogs(ctx context.Context, logs []*domain.LedgerLog) error {
	db := readDBFromCtx(ctx)

	var dbLedgerLogs []*ledgerLog
	for _, ledgerLog := range logs {
		dbLedgerLogs = append(dbLedgerLogs, marshallLedgerLog(ledgerLog))
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.CreateInBatches(dbLedgerLogs, 100).Error
	}); err != nil {
		return errors.Wrap(err, "failed to create ledger logs")
	}
	return nil
}

func (r LedgerPostgresRepository) ReadLedgerById(ctx context.Context, id uuid.UUID) (query.Ledger, error) {
	db := readDBFromCtx(ctx)

	dbLedger := ledger{}
	if err := db.First(&dbLedger, "id = ?", id).Error; err != nil {
		return query.Ledger{}, errors.Wrap(err, "find ledger by id failed")
	}

	return unmarshallLedgerToQuery(&dbLedger), nil
}

func (r LedgerPostgresRepository) ReadAllLedgersByAccountingPeriod(ctx context.Context, periodId uuid.UUID) ([]query.Ledger, error) {
	db := readDBFromCtx(ctx)

	var dbLedgers []ledger
	if err := db.Where("period_id = ?", periodId).Find(&dbLedgers).Error; err != nil {
		return nil, errors.Wrapf(err, "find ledgers by period %s failed", periodId)
	}

	queryLedgers := make([]query.Ledger, len(dbLedgers))
	for i, dbLedger := range dbLedgers {
		queryLedgers[i] = unmarshallLedgerToQuery(&dbLedger)
	}
	return queryLedgers, nil
}

func (r LedgerPostgresRepository) ReadAccountingPeriodById(ctx context.Context, id uuid.UUID) (query.AccountingPeriod, error) {
	db := readDBFromCtx(ctx)

	dbPeriod := accountingPeriod{}
	if err := db.First(&dbPeriod, "id = ?", id).Error; err != nil {
		return query.AccountingPeriod{}, errors.Wrap(err, "find accounting period by id failed")
	}

	return unmarshallPeriodToQuery(&dbPeriod), nil
}

func (r LedgerPostgresRepository) ReadLedgerLogsByAccountIdsAndTimes(ctx context.Context, accountIds []uuid.UUID, openingTime, endingTime time.Time) ([]query.LedgerLog, error) {
	db := readDBFromCtx(ctx)

	var dbLedgerLogs []ledgerLog
	if err := db.Where("account_id IN ? AND transaction_time >= ? AND transaction_time < ?", accountIds, openingTime, endingTime).Find(&dbLedgerLogs).Error; err != nil {
		return nil, errors.Wrapf(err, "find ledger logs by account and period failed")
	}

	var queryLedgerLogs []query.LedgerLog
	for _, dbLedgerLog := range dbLedgerLogs {
		queryLedgerLogs = append(queryLedgerLogs, unmarshallLedgerLogToQuery(&dbLedgerLog))
	}
	return queryLedgerLogs, nil
}

func (r LedgerPostgresRepository) ReadAllAccountingPeriods(ctx context.Context, sobId uuid.UUID) ([]query.AccountingPeriod, error) {
	db := readDBFromCtx(ctx)

	var dbPeriods []accountingPeriod
	if err := db.Where("sob_id = ?", sobId).Find(&dbPeriods).Error; err != nil {
		return nil, errors.Wrapf(err, "find accounting periods by sob %s failed", sobId)
	}

	var queryPeriods []query.AccountingPeriod
	for _, dbPeriod := range dbPeriods {
		queryPeriods = append(queryPeriods, unmarshallPeriodToQuery(&dbPeriod))
	}
	return queryPeriods, nil
}

func (r LedgerPostgresRepository) ReadOpenAccountingPeriod(ctx context.Context, sobId uuid.UUID) (query.AccountingPeriod, error) {
	db := readDBFromCtx(ctx)

	var dbPeriods []accountingPeriod
	if err := db.Where("sob_id = ? AND is_closed = false", sobId).Find(&dbPeriods).Error; err != nil {
		return query.AccountingPeriod{}, errors.Wrapf(err, "find open accounting period by sob %s failed", sobId)
	}

	if len(dbPeriods) != 1 {
		return query.AccountingPeriod{}, errors.Errorf("expects 1 open accounting period, but find %d", len(dbPeriods))
	}

	return unmarshallPeriodToQuery(&dbPeriods[0]), nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
