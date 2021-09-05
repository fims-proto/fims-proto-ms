package db

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain"

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

func (r AccountPostgresRepository) AddAccount(ctx context.Context, a *domain.Account) error {
	panic("not implemented") // TODO: Implement
}

func (r AccountPostgresRepository) Dataload(ctx context.Context, domainAccs []*domain.Account) error {
	if len(domainAccs) == 0 {
		return errors.New("empty account list")
	}

	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		// delete all within sob
		if err := tx.Where("sob_id = ?", domainAccs[0].Sob()).Delete(&account{}).Error; err != nil {
			return errors.Wrap(err, "accounts deletion failed")
		}

		// create all
		dbAccounts := []account{}
		for _, domainAcc := range domainAccs {
			dbAccounts = append(dbAccounts, *marshall(domainAcc))
		}
		if err := tx.Create(&dbAccounts).Error; err != nil {
			return errors.Wrap(err, "accounts create failed")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "account dataload failed")
	}

	return nil
}

func (r AccountPostgresRepository) ReadAllAccounts(ctx context.Context, sob string) ([]query.Account, error) {
	panic("not implemented") // TODO: Implement
}

func (r AccountPostgresRepository) ReadByNumber(ctx context.Context, sob string, accountNumber string) (query.Account, error) {
	db := readDBFromCtx(ctx)

	qas, err := r.readAccountWithSuperiorAccount(db, sob, accountNumber)
	if err != nil {
		return query.Account{}, errors.Wrapf(err, "failed to read account by number %s", accountNumber)
	}
	return qas, nil
}

func (r AccountPostgresRepository) readAccountWithSuperiorAccount(db *gorm.DB, sob, accountNumber string) (query.Account, error) {
	dbAccount := account{}
	if err := db.Where("sob_id = ? AND number = ?", sob, accountNumber).Find(&dbAccount).Error; err != nil {
		return query.Account{}, errors.Wrap(err, "read account by number failed")
	}

	result := unmarshallToQuery(&dbAccount)
	if dbAccount.SuperiorNumber == "" {
		return result, nil
	}
	superiorAccount, err := r.readAccountWithSuperiorAccount(db, sob, dbAccount.SuperiorNumber)
	if err != nil {
		return query.Account{}, err
	}
	result.SuperiorAccount = &superiorAccount
	return result, nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
