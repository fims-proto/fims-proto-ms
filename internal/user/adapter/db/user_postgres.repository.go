package db

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/database"

	"github/fims-proto/fims-proto-ms/internal/user/domain/user"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/user/app/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserPostgresRepository struct{}

func NewUserPostgresRepository() *UserPostgresRepository {
	return &UserPostgresRepository{}
}

func (r UserPostgresRepository) Migrate(ctx context.Context) error {
	db := database.ReadDBFromContext(ctx)

	if err := db.AutoMigrate(&userPO{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r UserPostgresRepository) UpsertUser(ctx context.Context, userId uuid.UUID, updateFn func(*user.User) (*user.User, error)) error {
	db := database.ReadDBFromContext(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		po := userPO{}
		var bo *user.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po, "id = ?", userId).Error; err != nil {
			bo, err = user.New(userId, nil)
			if err != nil {
				return errors.Wrap(err, "failed to create user")
			}
		} else {
			bo, err = userPOToBO(po)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal user")
			}
		}

		updatedBO, err := updateFn(bo)
		if err != nil {
			return errors.Wrap(err, "failed to update user in transaction")
		}

		po, err = userBOToPO(*updatedBO)
		if err != nil {
			return errors.Wrap(err, "failed to userBOToPO user")
		}
		return tx.Save(&po).Error
	})
}

// queries

func (r UserPostgresRepository) UserById(ctx context.Context, id uuid.UUID) (query.User, error) {
	db := database.ReadDBFromContext(ctx)

	po := userPO{}
	if err := db.Where("id = ?", id).First(&po).Error; err != nil {
		return query.User{}, errors.Wrapf(err, "failed to read id %s", id)
	}

	queryUser, err := userPOToDTO(po)
	if err != nil {
		return query.User{}, errors.Wrap(err, "failed to unmarshal user")
	}

	return queryUser, nil
}

func (r UserPostgresRepository) UsersByIds(ctx context.Context, ids []uuid.UUID) ([]query.User, error) {
	db := database.ReadDBFromContext(ctx)

	if len(ids) == 0 {
		return nil, nil
	}

	var userPOs []userPO
	if err := db.Where("id IN ?", ids).Find(&userPOs).Error; err != nil {
		return nil, errors.Wrapf(err, "failed to read ids")
	}

	var userDTOs []query.User
	for _, po := range userPOs {
		dto, err := userPOToDTO(po)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal user")
		}

		userDTOs = append(userDTOs, dto)
	}

	return userDTOs, nil
}
