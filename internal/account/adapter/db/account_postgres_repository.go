package db

import (
	"context"
	"strings"
	"time"

	"github/fims-proto/fims-proto-ms/internal/account/domain/ledger"

	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"github/fims-proto/fims-proto-ms/internal/account/domain/period"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"

	"gorm.io/gorm"
)

type AccountPostgresRepository struct{}

func NewAccountPostgresRepository() *AccountPostgresRepository {
	return &AccountPostgresRepository{}
}

func (r AccountPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	return db.AutoMigrate(&accountPO{}, &periodPO{}, &ledgerPO{})
}

func (r AccountPostgresRepository) InitialAccounts(ctx context.Context, accounts []*account.Account) error {
	if len(accounts) == 0 {
		return errors.New("empty Account configuration list")
	}

	db := readDBFromCtx(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		// delete all within sob
		if err := tx.Where("sob_id = ?", accounts[0].SobId()).Delete(&accountPO{}).Error; err != nil {
			return errors.Wrap(err, "accounts deletion failed")
		}

		// create all
		var accountPOs []accountPO
		for _, accountBO := range accounts {
			po, err := accountBOToPO(*accountBO)
			if err != nil {
				return errors.Wrap(err, "failed to map Account from BO to PO")
			}
			accountPOs = append(accountPOs, po)
		}
		return tx.CreateInBatches(&accountPOs, 100).Error
	})
}

func (r AccountPostgresRepository) CreatePeriod(ctx context.Context, period *period.Period, txFn func() error) error {
	db := readDBFromCtx(ctx)

	po := periodBOToPO(*period)

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&po).Error; err != nil {
			return err
		}

		// other actions in current transaction
		return txFn()
	})
}

func (r AccountPostgresRepository) CreateLedgers(ctx context.Context, ledgers []*ledger.Ledger) error {
	db := readDBFromCtx(ctx)

	var ledgerPOs []ledgerPO
	for _, bo := range ledgers {
		ledgerPOs = append(ledgerPOs, ledgerBOToPO(*bo))
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Omit("Account").CreateInBatches(&ledgerPOs, 500).Error
	})
}

func (r AccountPostgresRepository) UpdateLedgersByPeriodAndAccountIds(ctx context.Context, periodId uuid.UUID, accountIds []uuid.UUID, updateFn func(ledgers []*ledger.Ledger) ([]*ledger.Ledger, error)) error {
	db := readDBFromCtx(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		var ledgerPOs []ledgerPO
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("period_id = ? AND account_id IN ?", periodId, accountIds).
			Preload("Account").
			Find(&ledgerPOs).Error; err != nil {
			return err
		}

		var ledgerBOs []*ledger.Ledger
		for _, po := range ledgerPOs {
			bo, err := ledgerPOToBO(po)
			if err != nil {
				return errors.Wrap(err, "failed to map ledger")
			}
			ledgerBOs = append(ledgerBOs, bo)
		}

		updatedLedgers, err := updateFn(ledgerBOs)
		if err != nil {
			return errors.Wrap(err, "failed to update ledgers in transaction")
		}

		var updatedPOs []ledgerPO
		for _, updatedLedger := range updatedLedgers {
			updatedPOs = append(updatedPOs, ledgerBOToPO(*updatedLedger))
		}

		return tx.Save(&updatedPOs).Error
	})
}

// queries

func (r AccountPostgresRepository) PagingLedgersByPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID, pageable data.Pageable) (data.Page[query.Ledger], error) {
	db := readDBFromCtx(ctx)

	var ledgerPOS []ledgerPO

	db = db.Scopes(data.Filtering(pageable)).Where("sob_id = ? AND period_id = ?", sobId, periodId)

	var count int64
	if err := db.Model(&ledgerPO{}).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count ledgers")
	}

	if err := db.Scopes(data.Paging(pageable)).Preload("Account").Find(&ledgerPOS).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find ledgers by sobId %s and periodId %s", sobId, periodId)
	}

	var ledgerDTOs []query.Ledger
	for _, po := range ledgerPOS {
		dto, err := ledgerPOToDTO(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map ledger")
		}
		ledgerDTOs = append(ledgerDTOs, dto)
	}

	return data.NewPage(ledgerDTOs, pageable, int(count))
}

func (r AccountPostgresRepository) LedgersInPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID) ([]query.Ledger, error) {
	ledgers, err := r.PagingLedgersByPeriod(ctx, sobId, periodId, data.Unpaged())
	if err != nil {
		return nil, err
	}
	return ledgers.Content(), nil
}

func (r AccountPostgresRepository) PagingAccounts(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[query.Account], error) {
	db := readDBFromCtx(ctx)

	var accountPOS []accountPO

	db = db.Scopes(data.Filtering(pageable)).Where("sob_id = ?", sobId)

	var count int64
	if err := db.Model(&accountPO{}).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count accounts")
	}

	if err := db.Scopes(data.Paging(pageable)).Find(&accountPOS).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find accounts by sobId %s", sobId)
	}

	var accountDTOs []query.Account
	for _, po := range accountPOS {
		dto, err := accountPOToDTO(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map account")
		}
		accountDTOs = append(accountDTOs, dto)
	}

	return data.NewPage(accountDTOs, pageable, int(count))
}

func (r AccountPostgresRepository) AllAccounts(ctx context.Context, sobId uuid.UUID) ([]query.Account, error) {
	accounts, err := r.PagingAccounts(ctx, sobId, data.Unpaged())
	if err != nil {
		return nil, err
	}
	return accounts.Content(), nil
}

func (r AccountPostgresRepository) SuperiorAccounts(ctx context.Context, accountId uuid.UUID) ([]query.Account, error) {
	rawSql := `WITH RECURSIVE res AS (
		   SELECT *
		   FROM a_accounts
		   WHERE account_id = ?
		   UNION
		   SELECT a_accounts.*
		   FROM res
		   JOIN a_accounts ON a_accounts.account_id = res.superior_account_id
		)
		SELECT *
		FROM res
		WHERE account_id != ?`
	rawSql = strings.Join(strings.Fields(rawSql), " ") // normalize whitespaces

	db := readDBFromCtx(ctx)

	var accountPOs []accountPO
	if err := db.Raw(rawSql, accountId, accountId).Scan(&accountPOs).Error; err != nil {
		return nil, err
	}

	var accountDTOs []query.Account
	for _, po := range accountPOs {
		dto, err := accountPOToDTO(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map account")
		}

		accountDTOs = append(accountDTOs, dto)
	}

	return accountDTOs, nil
}

func (r AccountPostgresRepository) AccountsByIds(ctx context.Context, accountIds []uuid.UUID) ([]query.Account, error) {
	db := readDBFromCtx(ctx)

	var accountPOs []accountPO
	if err := db.Find(&accountPOs, "account_id IN ?", accountIds).Error; err != nil {
		return nil, errors.Wrap(err, "failed to read account")
	}

	var accountDTOs []query.Account
	for _, po := range accountPOs {
		dto, err := accountPOToDTO(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map Account configuration")
		}

		accountDTOs = append(accountDTOs, dto)
	}

	return accountDTOs, nil
}

func (r AccountPostgresRepository) AccountsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]query.Account, error) {
	db := readDBFromCtx(ctx)

	var accountPOs []accountPO
	if err := db.Find(&accountPOs, "sob_id = ? AND account_number IN ?", sobId, accountNumbers).Error; err != nil {
		return nil, errors.Wrap(err, "failed to read account")
	}

	var accountDTOs []query.Account
	for _, po := range accountPOs {
		dto, err := accountPOToDTO(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map Account configuration")
		}

		accountDTOs = append(accountDTOs, dto)
	}

	return accountDTOs, nil
}

func (r AccountPostgresRepository) PagingPeriods(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[query.Period], error) {
	db := readDBFromCtx(ctx)

	var periodPOs []periodPO

	db = db.Scopes(data.Filtering(pageable)).Where("sob_id = ?", sobId)

	var count int64
	if err := db.Model(&periodPO{}).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count periods")
	}

	if err := db.Scopes(data.Paging(pageable)).Find(&periodPOs).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find periods by sobId %s", sobId)
	}

	var periodDTOs []query.Period
	for _, po := range periodPOs {
		periodDTOs = append(periodDTOs, periodPOToDTO(po))
	}

	return data.NewPage(periodDTOs, pageable, int(count))
}

func (r AccountPostgresRepository) PeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (query.Period, error) {
	db := readDBFromCtx(ctx)

	var periodPOs []periodPO
	if err := db.
		Where("sob_id = ? AND opening_time <= ? AND (ending_time > ? OR ending_time = ?)", sobId, timePoint, timePoint, time.Time{}).
		Find(&periodPOs).Error; err != nil {
		return query.Period{}, errors.Wrap(err, "find period by id failed")
	}

	if len(periodPOs) != 1 {
		return query.Period{}, errors.Errorf("expected 1 but %d periods found", len(periodPOs))
	}

	return periodPOToDTO(periodPOs[0]), nil
}

func (r AccountPostgresRepository) PeriodById(ctx context.Context, periodId uuid.UUID) (query.Period, error) {
	db := readDBFromCtx(ctx)

	po := periodPO{}
	if err := db.First(&po, "id = ?", periodId).Error; err != nil {
		return query.Period{}, errors.Wrap(err, "failed to find period by id")
	}

	return periodPOToDTO(po), nil
}

func (r AccountPostgresRepository) PeriodsByIds(ctx context.Context, periodIds []uuid.UUID) ([]query.Period, error) {
	db := readDBFromCtx(ctx)

	var periodPOs []periodPO
	if err := db.Find(&periodPOs, "period_id IN ?", periodIds).Error; err != nil {
		return nil, errors.Wrap(err, "failed to read periods")
	}

	var periodDTOs []query.Period
	for _, po := range periodPOs {
		periodDTOs = append(periodDTOs, periodPOToDTO(po))
	}

	return periodDTOs, nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
