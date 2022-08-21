package db

import (
	"context"
	"strings"
	"time"

	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account_configuration"
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

	return db.AutoMigrate(&accountConfigurationPO{}, &periodPO{}, &accountPO{})
}

func (r AccountPostgresRepository) InitialAccountConfiguration(ctx context.Context, accountConfigurations []*account_configuration.AccountConfiguration) error {
	if len(accountConfigurations) == 0 {
		return errors.New("empty account configuration list")
	}

	db := readDBFromCtx(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		// delete all within sob
		if err := tx.Where("sob_id = ?", accountConfigurations[0].SobId()).Delete(&accountConfigurationPO{}).Error; err != nil {
			return errors.Wrap(err, "accounts deletion failed")
		}

		// create all
		var accountConfigurationPOs []accountConfigurationPO
		for _, accountConfigurationBO := range accountConfigurations {
			po, err := accountConfigurationBOToPO(*accountConfigurationBO)
			if err != nil {
				return errors.Wrap(err, "failed to map account from BO to PO")
			}
			accountConfigurationPOs = append(accountConfigurationPOs, po)
		}
		return tx.CreateInBatches(&accountConfigurationPOs, 100).Error
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

func (r AccountPostgresRepository) CreateAccounts(ctx context.Context, accounts []*account.Account) error {
	db := readDBFromCtx(ctx)

	var accountPOs []accountPO
	for _, bo := range accounts {
		accountPOs = append(accountPOs, accountBOToPO(*bo))
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Omit("Configuration", "Period").CreateInBatches(&accountPOs, 100).Error
	})
}

func (r AccountPostgresRepository) UpdateAccountsByPeriodAndIds(ctx context.Context, periodId uuid.UUID, accountIds []uuid.UUID, updateFn func(accounts []*account.Account) ([]*account.Account, error)) error {
	db := readDBFromCtx(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		var accountPOs []accountPO
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("period_id = ? AND account_id IN ?", periodId, accountIds).
			Preload("Configuration").
			Find(&accountPOs).Error; err != nil {
			return err
		}

		var accountBOs []*account.Account
		for _, vo := range accountPOs {
			bo, err := accountPOToBO(vo)
			if err != nil {
				return errors.Wrap(err, "failed to map account")
			}
			accountBOs = append(accountBOs, bo)
		}

		updatedAccounts, err := updateFn(accountBOs)
		if err != nil {
			return errors.Wrap(err, "failed to update account in transaction")
		}

		var updatedPOs []accountPO
		for _, updatedAccount := range updatedAccounts {
			updatedPOs = append(updatedPOs, accountBOToPO(*updatedAccount))
		}

		return tx.Save(&updatedPOs).Error
	})
}

// queries

func (r AccountPostgresRepository) PagingAccountsByPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID, pageable data.Pageable) (data.Page[query.Account], error) {
	db := readDBFromCtx(ctx)

	var accountPOs []accountPO

	db = db.Scopes(data.Filtering(pageable)).Where("sob_id = ? AND period_id = ?", sobId, periodId)

	var count int64
	if err := db.Model(&accountPO{}).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count accounts")
	}

	if err := db.Scopes(data.Paging(pageable)).Preload("Configuration").Find(&accountPOs).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find accounts by sobId %s and periodId %s", sobId, periodId)
	}

	var accountDTOs []query.Account
	for _, vo := range accountPOs {
		dto, err := accountPOToDTO(vo)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map account")
		}
		accountDTOs = append(accountDTOs, dto)
	}

	return data.NewPage(accountDTOs, pageable, int(count))
}

func (r AccountPostgresRepository) AccountsInPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID) ([]query.Account, error) {
	accounts, err := r.PagingAccountsByPeriod(ctx, sobId, periodId, data.Unpaged())
	if err != nil {
		return nil, err
	}
	return accounts.Content(), nil
}

func (r AccountPostgresRepository) PagingAccountConfigurations(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[query.AccountConfiguration], error) {
	db := readDBFromCtx(ctx)

	var accountConfigurationPOs []accountConfigurationPO

	db = db.Scopes(data.Filtering(pageable)).Where("sob_id = ?", sobId)

	var count int64
	if err := db.Model(&accountConfigurationPO{}).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count account configurations")
	}

	if err := db.Scopes(data.Paging(pageable)).Find(&accountConfigurationPOs).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to find account configurations by sobId %s", sobId)
	}

	var configDTOs []query.AccountConfiguration
	for _, po := range accountConfigurationPOs {
		dto, err := accountConfigurationPOToDTO(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map account configuration")
		}
		configDTOs = append(configDTOs, dto)
	}

	return data.NewPage(configDTOs, pageable, int(count))
}

func (r AccountPostgresRepository) AllAccountConfigurations(ctx context.Context, sobId uuid.UUID) ([]query.AccountConfiguration, error) {
	configuration, err := r.PagingAccountConfigurations(ctx, sobId, data.Unpaged())
	if err != nil {
		return nil, err
	}
	return configuration.Content(), nil
}

func (r AccountPostgresRepository) SuperiorAccountConfigurations(ctx context.Context, accountId uuid.UUID) ([]query.AccountConfiguration, error) {
	rawSql := `WITH RECURSIVE res AS (
		   SELECT *
		   FROM account_configurations
		   WHERE account_id = ?
		   UNION
		   SELECT account_configurations.*
		   FROM res
		   JOIN account_configurations ON account_configurations.account_id = res.superior_account_id
		)
		SELECT *
		FROM res
		WHERE account_id != ?`
	rawSql = strings.Join(strings.Fields(rawSql), " ") // normalize whitespaces

	db := readDBFromCtx(ctx)

	var accountConfigurations []accountConfigurationPO
	if err := db.Raw(rawSql, accountId, accountId).Scan(&accountConfigurations).Error; err != nil {
		return nil, err
	}

	var acDTOs []query.AccountConfiguration
	for _, po := range accountConfigurations {
		dto, err := accountConfigurationPOToDTO(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map account configuration")
		}

		acDTOs = append(acDTOs, dto)
	}

	return acDTOs, nil
}

func (r AccountPostgresRepository) AccountConfigurationsByIds(ctx context.Context, accountIds []uuid.UUID) ([]query.AccountConfiguration, error) {
	db := readDBFromCtx(ctx)

	var configPOs []accountConfigurationPO
	if err := db.Find(&configPOs, "WHERE account_id IN ?", accountIds).Error; err != nil {
		return nil, errors.Wrap(err, "failed to read account configuration")
	}

	var configDTOs []query.AccountConfiguration
	for _, po := range configPOs {
		dto, err := accountConfigurationPOToDTO(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map account configuration")
		}

		configDTOs = append(configDTOs, dto)
	}

	return configDTOs, nil
}

func (r AccountPostgresRepository) AccountConfigurationsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]query.AccountConfiguration, error) {
	db := readDBFromCtx(ctx)

	var configPOs []accountConfigurationPO
	if err := db.Find(&configPOs, "WHERE sob_id = ? AND account_number IN ?", sobId, accountNumbers).Error; err != nil {
		return nil, errors.Wrap(err, "failed to read account configuration")
	}

	var configDTOs []query.AccountConfiguration
	for _, po := range configPOs {
		dto, err := accountConfigurationPOToDTO(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to map account configuration")
		}

		configDTOs = append(configDTOs, dto)
	}

	return configDTOs, nil
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
	if err := db.Find(&periodPOs, "WHERE period_id IN ?", periodIds).Error; err != nil {
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
