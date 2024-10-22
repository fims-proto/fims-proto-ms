package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	"github/fims-proto/fims-proto-ms/internal/user/app/query"
	"github/fims-proto/fims-proto-ms/internal/user/domain/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserPostgresRepository struct {
	dataSource datasource.DataSource
}

func NewUserPostgresRepository(dataSource datasource.DataSource) *UserPostgresRepository {
	return &UserPostgresRepository{
		dataSource: dataSource,
	}
}

func (r UserPostgresRepository) Migrate(ctx context.Context) error {
	db := r.dataSource.GetConnection(ctx)

	if err := db.AutoMigrate(&userPO{}); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	return nil
}

func (r UserPostgresRepository) UpsertUser(ctx context.Context, userId uuid.UUID, updateFn func(*user.User) (*user.User, error)) error {
	db := r.dataSource.GetConnection(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		po := userPO{Id: userId}
		var bo *user.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po).Error; err != nil {
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

		po = userBOToPO(*updatedBO)
		return tx.Save(&po).Error
	})
}

// queries

func (r UserPostgresRepository) UserById(ctx context.Context, id uuid.UUID) (query.User, error) {
	db := r.dataSource.GetConnection(ctx)

	po := userPO{Id: id}
	if err := db.First(&po).Error; err != nil {
		return query.User{}, fmt.Errorf("failed to read id %s: %w", id, err)
	}

	return userPOToDTO(po), nil
}

func (r UserPostgresRepository) UsersByIds(ctx context.Context, ids []uuid.UUID) ([]query.User, error) {
	db := r.dataSource.GetConnection(ctx)

	if len(ids) == 0 {
		return nil, nil
	}

	var userPOs []userPO
	if err := db.Where("id IN ?", ids).Find(&userPOs).Error; err != nil {
		return nil, fmt.Errorf("failed to read users: %w", err)
	}

	return converter.POsToDTOs(userPOs, userPOToDTO), nil
}
