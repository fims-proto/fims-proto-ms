package db

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LedgerPostgresRepository struct{}

func NewLedgerPostgresRepository() *LedgerPostgresRepository {
	return &LedgerPostgresRepository{}
}

func (r LedgerPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	if err := db.AutoMigrate(&ledger{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r LedgerPostgresRepository) AddLedger(ctx context.Context, l *domain.Ledger) error {
	db := readDBFromCtx(ctx)

	dbLedger := marshall(l)

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(dbLedger).Error
	}); err != nil {
		return errors.Wrap(err, "failed to create ledger")
	}
	return nil
}

func (r LedgerPostgresRepository) UpdateLedgers(ctx context.Context, sob string, ledgerNumbers []string, updateFn func(ledgers []*domain.Ledger) ([]*domain.Ledger, error)) error {
	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		dbLedgers := []ledger{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("sob_id = ? AND number IN ?", sob, ledgerNumbers).Find(&dbLedgers).Error; err != nil {
			return errors.Wrap(err, "find ledgers failed")
		}

		if len(dbLedgers) == 0 {
			return errors.Errorf("no ledger found by sob %s and numbers %s", sob, ledgerNumbers)
		}

		domainLedgers := []*domain.Ledger{}
		for _, dbLedger := range dbLedgers {
			domainLedger, err := unmarshallToDomain(&dbLedger)
			if err != nil {
				return errors.Wrapf(err, "failed to unmarshall ledger %s", dbLedger.Number)
			}
			domainLedgers = append(domainLedgers, domainLedger)
		}

		updatedDomainLedgers, err := updateFn(domainLedgers)
		if err != nil {
			return errors.Wrap(err, "failed to update ledgers in transaction")
		}

		dbLedgers = []ledger{}
		for _, updatedDomainLedger := range updatedDomainLedgers {
			dbLedgers = append(dbLedgers, *marshall(updatedDomainLedger))
		}

		if err := tx.Save(&dbLedgers).Error; err != nil {
			return errors.Wrap(err, "failed to save ledgers")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "failed to updated ledgers")
	}

	return nil
}

func (r LedgerPostgresRepository) Dataload(ctx context.Context, domainLedgers []*domain.Ledger) error {
	if len(domainLedgers) == 0 {
		return errors.New("empty ledger list")
	}

	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		// delete all within sob
		if err := tx.Where("sob_id = ?", domainLedgers[0].Sob()).Delete(&ledger{}).Error; err != nil {
			return errors.Wrap(err, "ledgers deletion failed")
		}

		// create all
		dbLedgers := []ledger{}
		for _, domainLedger := range domainLedgers {
			dbLedgers = append(dbLedgers, *marshall(domainLedger))
		}
		if err := tx.Create(&dbLedgers).Error; err != nil {
			return errors.Wrap(err, "ledgers create failed")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "ledger dataload failed")
	}

	return nil
}

func (r LedgerPostgresRepository) ReadAllLedgers(ctx context.Context, sob string) ([]query.Ledger, error) {
	db := readDBFromCtx(ctx)

	dbLedgers := []ledger{}
	if err := db.Where("sob_id = ?", sob).Find(&dbLedgers).Error; err != nil {
		return []query.Ledger{}, errors.Wrapf(err, "find ledgers by sob %s failed", sob)
	}

	queryLedgers := []query.Ledger{}
	for _, dbLedger := range dbLedgers {
		queryLedgers = append(queryLedgers, unmarshallToQuery(&dbLedger))
	}
	return queryLedgers, nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
