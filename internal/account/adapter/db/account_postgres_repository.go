package db

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain"

	"github.com/google/uuid"
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

func (r AccountPostgresRepository) CreateAccount(context.Context, *domain.Account) error {
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
			dbAccount, err := marshal(domainAcc)
			if err != nil {
				return errors.Wrap(err, "failed to marshal account")
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

func (r AccountPostgresRepository) ReadAccounts(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[query.Account], error) {
	db := readDBFromCtx(ctx)

	var dbAccounts []account

	db = data.AddFilter(pageable, db).Where("sob_id = ?", sobId)

	var count int64
	if err := db.Model(&account{}).Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, "count accounts failed")
	}

	db = data.AddPaging(pageable, db)

	if err := db.Find(&dbAccounts).Error; err != nil {
		return nil, errors.Wrapf(err, "find accounts by sobId %s failed", sobId)
	}

	var queryAccounts []query.Account
	for _, dbAccount := range dbAccounts {
		queryAccount, err := unmarshalToQuery(dbAccount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal account")
		}
		queryAccounts = append(queryAccounts, queryAccount)
	}

	accountsPage, err := data.NewPage(queryAccounts, pageable, int(count))
	if err != nil {
		return nil, errors.Wrap(err, "wrap to page failed")
	}

	return accountsPage, nil
}

func (r AccountPostgresRepository) ReadAccountsWithSuperiorsByIds(ctx context.Context, accountIds []uuid.UUID) ([]query.Account, error) {
	db := readDBFromCtx(ctx)

	queryAccounts, err := r.readAccountsWithSuperiors(db, accountIds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read accounts by ids")
	}
	return queryAccounts, nil
}

func (r AccountPostgresRepository) ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]query.Account, error) {
	db := readDBFromCtx(ctx)

	if len(accountIds) == 0 {
		return nil, nil
	}

	var dbAccounts []account
	if err := db.Where("id IN ?", accountIds).Find(&dbAccounts).Error; err != nil {
		return nil, errors.Wrapf(err, "find accounts by IDs")
	}

	queryAccounts := make(map[uuid.UUID]query.Account)
	for _, dbAccount := range dbAccounts {
		queryAccount, err := unmarshalToQuery(dbAccount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal account")
		}
		queryAccounts[dbAccount.Id] = queryAccount
	}
	return queryAccounts, nil
}

func (r AccountPostgresRepository) ReadAccountByNumber(ctx context.Context, sobId uuid.UUID, accountNumber string) (query.Account, error) {
	db := readDBFromCtx(ctx)

	dbAccount := account{}
	if err := db.Where("sob_id = ? AND account_number = ?", sobId, accountNumber).
		Find(&dbAccount).Error; err != nil {
		return query.Account{}, errors.Wrap(err, "read account by number and superior number failed")
	}
	result, err := unmarshalToQuery(dbAccount)
	if err != nil {
		return query.Account{}, errors.Wrap(err, "failed to unmarshal account")
	}
	return result, nil
}

func (r AccountPostgresRepository) readAccountsWithSuperiors(db *gorm.DB, ids []uuid.UUID) ([]query.Account, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var dbAccounts []account
	if err := db.Where("id IN ?", ids).Find(&dbAccounts).Error; err != nil {
		return nil, errors.Wrap(err, "failed to read accounts by ids")
	}

	var queryAccounts []*query.Account
	var superiorIds []uuid.UUID
	superior2AccountMap := make(map[uuid.UUID]*query.Account)
	for _, dbAccount := range dbAccounts {
		queryAccount, err := unmarshalToQuery(dbAccount)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal db account")
		}
		queryAccounts = append(queryAccounts, &queryAccount)

		if dbAccount.SuperiorAccountId != uuid.Nil {
			superiorIds = append(superiorIds, dbAccount.SuperiorAccountId)
			superior2AccountMap[dbAccount.SuperiorAccountId] = &queryAccount
		}
	}

	superiorAccounts, err := r.readAccountsWithSuperiors(db, superiorIds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read superior accounts")
	}

	for _, superiorAccount := range superiorAccounts {
		childAccount, ok := superior2AccountMap[superiorAccount.Id]
		if !ok {
			return nil, errors.Errorf("failed to find child account by superior %s", superiorAccount.Id)
		}
		childAccount.SuperiorAccount = &superiorAccount
	}

	var res []query.Account
	for _, queryAccount := range queryAccounts {
		res = append(res, *queryAccount)
	}

	return res, nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
