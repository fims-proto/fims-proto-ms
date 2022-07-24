package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/user/app/query"
	"github/fims-proto/fims-proto-ms/internal/user/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserPostgresRepository struct{}

func NewUserPostgresRepository() *UserPostgresRepository {
	return &UserPostgresRepository{}
}

func (r UserPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	if err := db.AutoMigrate(&user{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r UserPostgresRepository) ReadById(ctx context.Context, id uuid.UUID) (query.User, error) {
	db := readDBFromCtx(ctx)

	dbUser := user{}
	if err := db.Where("id = ?", id).First(&dbUser).Error; err != nil {
		return query.User{}, errors.Wrapf(err, "failed to read id %s", id)
	}

	queryUser, err := unmarshalToQuery(dbUser)
	if err != nil {
		return query.User{}, errors.Wrap(err, "failed to unmarshal user")
	}
	return queryUser, nil
}

func (r UserPostgresRepository) ReadByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]query.User, error) {
	db := readDBFromCtx(ctx)

	if len(ids) == 0 {
		return nil, nil
	}

	var dbUsers []user
	if err := db.Where("id IN ?", ids).Find(&dbUsers).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to read ids")
	}

	queryUsers := make(map[uuid.UUID]query.User)
	for _, dbUser := range dbUsers {
		queryUser, err := unmarshalToQuery(dbUser)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal user")
		}
		queryUsers[dbUser.Id] = queryUser
	}

	return queryUsers, nil
}

func (r UserPostgresRepository) UpdateUser(ctx context.Context, id uuid.UUID, updateFn func(*domain.User) (*domain.User, error)) error {
	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		dbUser := user{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&dbUser, "id = ?", id).Error; err != nil {
			return errors.Wrap(err, "failed to find user")
		}

		domainUser, err := unmarshalToDomain(dbUser)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal user")
		}

		updatedDomainUser, err := updateFn(domainUser)
		if err != nil {
			return errors.Wrap(err, "failed to update user in transaction")
		}

		dbUser, err = marshal(*updatedDomainUser)
		if err != nil {
			return errors.Wrap(err, "failed to marshal user")
		}
		if err = tx.Save(dbUser).Error; err != nil {
			return errors.Wrap(err, "failed to save user")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "failed to update user")
	}
	return nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
