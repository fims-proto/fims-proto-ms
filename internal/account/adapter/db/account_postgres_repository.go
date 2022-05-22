package db

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

type AccountPostgresRepository struct{}

func NewAccountPostgresRepository() *AccountPostgresRepository {
	return &AccountPostgresRepository{}
}

func (r AccountPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	if err := db.AutoMigrate(&account{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r AccountPostgresRepository) CreateAccount(ctx context.Context, a *domain.Account) error {
	panic("not implemented") // TODO: Implement
}

func (r AccountPostgresRepository) DataLoad(ctx context.Context, domainAccounts []*domain.Account) error {
	if len(domainAccounts) == 0 {
		return errors.New("empty account list")
	}

	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		// delete all within sob
		if err := tx.Where("sob_id = ?", domainAccounts[0].SobId()).Delete(&account{}).Error; err != nil {
			return errors.Wrap(err, "accounts deletion failed")
		}

		// create all
		var dbAccounts []account
		for _, domainAcc := range domainAccounts {
			dbAccount, err := marshall(domainAcc)
			if err != nil {
				return errors.Wrap(err, "failed to marshall account")
			}
			dbAccounts = append(dbAccounts, *dbAccount)
		}
		if err := tx.CreateInBatches(&dbAccounts, 100).Error; err != nil {
			return errors.Wrap(err, "accounts create failed")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "account data load failed")
	}

	return nil
}

func (r AccountPostgresRepository) ReadAllAccounts(ctx context.Context, sobId uuid.UUID) ([]query.Account, error) {
	db := readDBFromCtx(ctx)

	var dbAccounts []account
	if err := db.Where("sob_id = ?", sobId).Find(&dbAccounts).Error; err != nil {
		return nil, errors.Wrapf(err, "find accounts by sobId %s failed", sobId)
	}

	var queryAccounts []query.Account
	for _, dbAccount := range dbAccounts {
		queryAccount, err := unmarshallToQuery(&dbAccount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshall account")
		}
		queryAccounts = append(queryAccounts, *queryAccount)
	}
	return queryAccounts, nil
}

func (r AccountPostgresRepository) ReadById(ctx context.Context, accountId uuid.UUID) (query.Account, error) {
	db := readDBFromCtx(ctx)

	qas, err := r.readAccountWithSuperiorAccount(db, accountId)
	if err != nil {
		return query.Account{}, errors.Wrapf(err, "failed to read account %s", accountId.String())
	}
	return qas, nil
}

func (r AccountPostgresRepository) ReadByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]*query.Account, error) {
	db := readDBFromCtx(ctx)

	if len(accountIds) == 0 {
		return nil, nil
	}

	var dbAccounts []account
	if err := db.Where("id IN ?", accountIds).Find(&dbAccounts).Error; err != nil {
		return nil, errors.Wrapf(err, "find accounts by IDs")
	}

	queryAccounts := make(map[uuid.UUID]*query.Account)
	for _, dbAccount := range dbAccounts {
		queryAccount, err := unmarshallToQuery(&dbAccount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshall account")
		}
		queryAccounts[dbAccount.Id] = queryAccount
	}
	return queryAccounts, nil
}

func (r AccountPostgresRepository) ReadByAccountNumber(ctx context.Context, sobId uuid.UUID, levelNumber int, superiorNumbers []int) (query.Account, error) {
	var superiorNumbersDb pgtype.Int4Array
	if err := superiorNumbersDb.Set(superiorNumbers); err != nil {
		return query.Account{}, errors.Wrap(err, "convert []int to Int4Array failed")
	}

	db := readDBFromCtx(ctx)

	dbAccount := account{}
	if err := db.Where("sob_id = ? AND level_number = ? AND superior_numbers = ?", sobId, levelNumber, superiorNumbersDb).
		Find(&dbAccount).Error; err != nil {
		return query.Account{}, errors.Wrap(err, "read account by number and superior number failed")
	}
	result, err := unmarshallToQuery(&dbAccount)
	if err != nil {
		return query.Account{}, errors.Wrap(err, "failed to unmarshall account")
	}
	return *result, nil
}

func (r AccountPostgresRepository) readAccountWithSuperiorAccount(db *gorm.DB, accountId uuid.UUID) (query.Account, error) {
	dbAccount := account{}
	if err := db.Where("id = ?", accountId).Find(&dbAccount).Error; err != nil {
		return query.Account{}, errors.Wrap(err, "read account by id failed")
	}

	result, err := unmarshallToQuery(&dbAccount)
	if err != nil {
		return query.Account{}, errors.Wrap(err, "failed to unmarshall account")
	}
	if dbAccount.SuperiorAccountId == uuid.Nil {
		return *result, nil
	}
	superiorAccount, err := r.readAccountWithSuperiorAccount(db, dbAccount.SuperiorAccountId)
	if err != nil {
		return query.Account{}, err
	}
	result.SuperiorAccount = &superiorAccount
	return *result, nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
