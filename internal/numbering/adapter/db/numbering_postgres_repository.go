package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/numbering/app/query"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NumberingPostgresRepository struct{}

func NewNumberingPostgresRepository() *NumberingPostgresRepository {
	return &NumberingPostgresRepository{}
}

func (r NumberingPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	if err := db.AutoMigrate(&identifierConfiguration{}, &identifier{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r NumberingPostgresRepository) CreateIdentifierConfiguration(ctx context.Context, configuration *domain.IdentifierConfiguration) error {
	db := readDBFromCtx(ctx)

	dbConfig, err := marshalIdentifierConfiguration(configuration)
	if err != nil {
		return err
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(dbConfig).Error
	}); err != nil {
		return errors.Wrap(err, "create identifier configuration failed")
	}

	return nil
}

func (r NumberingPostgresRepository) UpdateIdentifierConfiguration(ctx context.Context, id uuid.UUID, updateFn func(config *domain.IdentifierConfiguration) (*domain.IdentifierConfiguration, error)) error {
	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		dbConfig := &identifierConfiguration{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(dbConfig, "id = ?", id).Error; err != nil {
			return err
		}

		domainConfig, err := unmarshalToIdentConfigDomain(dbConfig)
		if err != nil {
			return errors.Wrap(err, "unmarshal identifier configuration failed")
		}

		updatedConfig, err := updateFn(domainConfig)
		if err != nil {
			return errors.Wrap(err, "update identifier configuration in transaction failed")
		}

		dbConfig, err = marshalIdentifierConfiguration(updatedConfig)
		if err != nil {
			return errors.Wrap(err, "marshal identifier configuration failed")
		}
		if err = tx.Save(dbConfig).Error; err != nil {
			return errors.Wrap(err, "save identifier configuration failed")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "update identifier configuration failed")
	}

	return nil
}

func (r NumberingPostgresRepository) CreateIdentifier(ctx context.Context, identifier *domain.Identifier) error {
	db := readDBFromCtx(ctx)

	dbIdent := marshalIdentifier(identifier)

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(dbIdent).Error
	}); err != nil {
		return errors.Wrap(err, "create identifier failed")
	}

	return nil
}

func (r NumberingPostgresRepository) ResolveIdentifierConfiguration(ctx context.Context, targetBusinessObject string, objectsToMatch map[string]string) (query.IdentifierConfiguration, error) {
	db := readDBFromCtx(ctx)

	var dbConfigs []identifierConfiguration
	if err := db.Where("target_business_object = ?", targetBusinessObject).Find(&dbConfigs).Error; err != nil {
		return query.IdentifierConfiguration{}, errors.Wrap(err, "failed to find identifier configuration by business object")
	}

	for _, dbConfig := range dbConfigs {
		domainConfig, err := unmarshalToIdentConfigDomain(&dbConfig)
		if err != nil {
			return query.IdentifierConfiguration{}, errors.Wrap(err, "unmarshal identifier configuration failed")
		}
		if domainConfig.IsMatchProperties(objectsToMatch) {
			queryConfig, _ := unmarshalToIdentConfigQuery(dbConfig)
			return queryConfig, nil
		}
	}

	return query.IdentifierConfiguration{}, errors.New("no identifier configuration find")
}

func (r NumberingPostgresRepository) IdentifierById(ctx context.Context, id uuid.UUID) (query.Identifier, error) {
	db := readDBFromCtx(ctx)

	dbIdentifier := &identifier{}
	if err := db.First(dbIdentifier, "id = ?", id).Error; err != nil {
		return query.Identifier{}, errors.Wrap(err, "failed to read identifier by id")
	}

	return unmarshalToIdentifier(*dbIdentifier), nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
