package db

import (
	"context"
	"errors"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	"github/fims-proto/fims-proto-ms/internal/numbering/app/query"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier_configuration"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NumberingPostgresRepository struct {
	dataSource datasource.DataSource
}

func NewNumberingPostgresRepository(dataSource datasource.DataSource) *NumberingPostgresRepository {
	return &NumberingPostgresRepository{
		dataSource: dataSource,
	}
}

func (r NumberingPostgresRepository) Migrate(ctx context.Context) error {
	db := r.dataSource.GetConnection(ctx)

	if err := db.AutoMigrate(&identifierConfigurationPO{}, &identifierPO{}); err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}
	return nil
}

func (r NumberingPostgresRepository) CreateIdentifierConfiguration(
	ctx context.Context,
	domainConfig *identifier_configuration.IdentifierConfiguration,
) error {
	db := r.dataSource.GetConnection(ctx)

	dbConfig, err := identifierConfigurationBOToPO(*domainConfig)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&dbConfig).Error
	})
}

func (r NumberingPostgresRepository) UpdateIdentifierConfiguration(
	ctx context.Context,
	id uuid.UUID,
	updateFn func(config *identifier_configuration.IdentifierConfiguration) (*identifier_configuration.IdentifierConfiguration, error),
) error {
	db := r.dataSource.GetConnection(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		po := identifierConfigurationPO{Id: id}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po).Error; err != nil {
			return err
		}

		bo, err := identifierConfigurationPOToBO(po)
		if err != nil {
			return err
		}

		updatedBO, err := updateFn(bo)
		if err != nil {
			return fmt.Errorf("failed to update identifier configuration: %w", err)
		}

		po, err = identifierConfigurationBOToPO(*updatedBO)
		if err != nil {
			return err
		}
		return tx.Save(&po).Error
	})
}

func (r NumberingPostgresRepository) CreateIdentifier(ctx context.Context, bo *identifier.Identifier) error {
	db := r.dataSource.GetConnection(ctx)

	po := identifierBOToPO(*bo)

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&po).Error
	})
}

func (r NumberingPostgresRepository) ResolveIdentifierConfiguration(
	ctx context.Context,
	targetBusinessObject string,
	objectsToMatch map[string]string,
) (query.IdentifierConfiguration, error) {
	db := r.dataSource.GetConnection(ctx)

	var configPOs []identifierConfigurationPO
	if err := db.Where("target_business_object = ?", targetBusinessObject).Find(&configPOs).Error; err != nil {
		return query.IdentifierConfiguration{}, fmt.Errorf("failed to find identifier configuration by business object: %w", err)
	}

	for _, po := range configPOs {
		bo, err := identifierConfigurationPOToBO(po)
		if err != nil {
			return query.IdentifierConfiguration{}, fmt.Errorf("failed to unmarshal identifier configuration: %w", err)
		}

		if bo.IsMatchProperties(objectsToMatch) {
			dto, _ := identifierConfigurationPOToDTO(po)
			return dto, nil
		}
	}

	return query.IdentifierConfiguration{}, errors.New("no identifier configuration find")
}

func (r NumberingPostgresRepository) IdentifierById(ctx context.Context, id uuid.UUID) (query.Identifier, error) {
	db := r.dataSource.GetConnection(ctx)

	po := identifierPO{Id: id}
	if err := db.First(&po).Error; err != nil {
		return query.Identifier{}, fmt.Errorf("failed to read identifier by id: %w", err)
	}

	return identifierPOToDTO(po), nil
}
