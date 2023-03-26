package db

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/database"

	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier_configuration"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/numbering/app/query"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NumberingPostgresRepository struct{}

func NewNumberingPostgresRepository() *NumberingPostgresRepository {
	return &NumberingPostgresRepository{}
}

func (r NumberingPostgresRepository) Migrate(ctx context.Context) error {
	db := database.ReadDBFromContext(ctx)

	if err := db.AutoMigrate(&identifierConfigurationPO{}, &identifierPO{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r NumberingPostgresRepository) CreateIdentifierConfiguration(ctx context.Context, domainConfig *identifier_configuration.IdentifierConfiguration) error {
	db := database.ReadDBFromContext(ctx)

	dbConfig, err := identifierConfigurationBOToPO(*domainConfig)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&dbConfig).Error
	})
}

func (r NumberingPostgresRepository) UpdateIdentifierConfiguration(ctx context.Context, id uuid.UUID, updateFn func(config *identifier_configuration.IdentifierConfiguration) (*identifier_configuration.IdentifierConfiguration, error)) error {
	db := database.ReadDBFromContext(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		po := identifierConfigurationPO{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&po, "id = ?", id).Error; err != nil {
			return err
		}

		bo, err := identifierConfigurationPOToBO(po)
		if err != nil {
			return errors.Wrap(err, "unmarshal identifier configuration failed")
		}

		updatedBO, err := updateFn(bo)
		if err != nil {
			return errors.Wrap(err, "update identifier configuration in transaction failed")
		}

		po, err = identifierConfigurationBOToPO(*updatedBO)
		if err != nil {
			return errors.Wrap(err, "marshal identifier configuration failed")
		}
		return tx.Save(&po).Error
	})
}

func (r NumberingPostgresRepository) CreateIdentifier(ctx context.Context, bo *identifier.Identifier) error {
	db := database.ReadDBFromContext(ctx)

	po := identifierBOToPO(*bo)

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&po).Error
	})
}

func (r NumberingPostgresRepository) ResolveIdentifierConfiguration(ctx context.Context, targetBusinessObject string, objectsToMatch map[string]string) (query.IdentifierConfiguration, error) {
	db := database.ReadDBFromContext(ctx)

	var configPOs []identifierConfigurationPO
	if err := db.Where("target_business_object = ?", targetBusinessObject).Find(&configPOs).Error; err != nil {
		return query.IdentifierConfiguration{}, errors.Wrap(err, "failed to find identifier configuration by business object")
	}

	for _, po := range configPOs {
		bo, err := identifierConfigurationPOToBO(po)
		if err != nil {
			return query.IdentifierConfiguration{}, errors.Wrap(err, "unmarshal identifier configuration failed")
		}

		if bo.IsMatchProperties(objectsToMatch) {
			dto, _ := identifierConfigurationPOToDTO(po)
			return dto, nil
		}
	}

	return query.IdentifierConfiguration{}, errors.New("no identifier configuration find")
}

func (r NumberingPostgresRepository) IdentifierById(ctx context.Context, id uuid.UUID) (query.Identifier, error) {
	db := database.ReadDBFromContext(ctx)

	po := identifierPO{}
	if err := db.First(&po, "id = ?", id).Error; err != nil {
		return query.Identifier{}, errors.Wrap(err, "failed to read identifier by id")
	}

	return identifierPOToDTO(po), nil
}
