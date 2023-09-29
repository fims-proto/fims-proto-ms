package db

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/database"

	"github/fims-proto/fims-proto-ms/internal/user/domain/user"

	"github.com/google/uuid"
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
		return fmt.Errorf("failed to migrate: %w", err)
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
				return fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			bo, err = userPOToBO(po)
			if err != nil {
				return fmt.Errorf("failed to unmarshal user: %w", err)
			}
		}

		updatedBO, err := updateFn(bo)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		po, err = userBOToPO(*updatedBO)
		if err != nil {
			return err
		}
		return tx.Save(&po).Error
	})
}

// queries

func (r UserPostgresRepository) UserById(ctx context.Context, id uuid.UUID) (query.User, error) {
	db := database.ReadDBFromContext(ctx)

	po := userPO{}
	if err := db.Where("id = ?", id).First(&po).Error; err != nil {
		return query.User{}, fmt.Errorf("failed to read id %s: %w", id, err)
	}

	queryUser, err := userPOToDTO(po)
	if err != nil {
		return query.User{}, err
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
		return nil, fmt.Errorf("failed to read users: %w", err)
	}

	var userDTOs []query.User
	for _, po := range userPOs {
		dto, err := userPOToDTO(po)
		if err != nil {
			return nil, err
		}
		userDTOs = append(userDTOs, dto)
	}

	return userDTOs, nil
}
