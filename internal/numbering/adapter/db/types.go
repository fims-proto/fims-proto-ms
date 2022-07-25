package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/numbering/app/query"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain"
)

type propertyMatcher struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type identifierConfiguration struct {
	Id                   uuid.UUID    `gorm:"type:uuid"`
	TargetBusinessObject string       `gorm:"uniqueIndex:identifierConfigs_target_matcher_key"`
	PropertyMatchers     pgtype.JSONB `gorm:"uniqueIndex:identifierConfigs_target_matcher_key"`
	Counter              int
	Prefix               string
	Suffix               string
	CreatedAt            time.Time `gorm:"<-:create"`
	UpdatedAt            time.Time
}

type identifier struct {
	Id                        uuid.UUID `gorm:"type:uuid"`
	IdentifierConfigurationId uuid.UUID `gorm:"type:uuid;uniqueIndex:identifiers_configuration_identifier_key"`
	Identifier                string    `gorm:"uniqueIndex:identifiers_configuration_identifier_key"`
	CreatedAt                 time.Time `gorm:"<-:create"`
}

func marshalPropertyMatchers(domainMatchers []domain.PropertyMatcher) (pgtype.JSONB, error) {
	var matchers []propertyMatcher
	for _, matcher := range domainMatchers {
		matchers = append(matchers, propertyMatcher{
			Name:  matcher.Name(),
			Value: matcher.Value(),
		})
	}
	var dbMatchers pgtype.JSONB
	if err := dbMatchers.Set(matchers); err != nil {
		return pgtype.JSONB{}, errors.Wrapf(err, "failed to convert %v to pgtype.JSONB", domainMatchers)
	}

	return dbMatchers, nil
}

func marshalIdentifierConfiguration(config domain.IdentifierConfiguration) (identifierConfiguration, error) {
	matchers, err := marshalPropertyMatchers(config.PropertyMatchers())
	if err != nil {
		return identifierConfiguration{}, err
	}

	return identifierConfiguration{
		Id:                   config.Id(),
		TargetBusinessObject: config.TargetBusinessObject(),
		PropertyMatchers:     matchers,
		Counter:              config.Counter(),
		Prefix:               config.Prefix(),
		Suffix:               config.Suffix(),
	}, nil
}

func marshalIdentifier(ident domain.Identifier) identifier {
	return identifier{
		Id:                        ident.Id(),
		IdentifierConfigurationId: ident.IdentifierConfigurationId(),
		Identifier:                ident.Identifier(),
	}
}

func unmarshalPropertyMatchers(dbMatchers pgtype.JSONB) ([]propertyMatcher, error) {
	var matchers []propertyMatcher
	if err := dbMatchers.AssignTo(&matchers); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal property matchers")
	}
	return matchers, nil
}

func unmarshalToIdentConfigDomain(dbConfig identifierConfiguration) (*domain.IdentifierConfiguration, error) {
	matchers, err := unmarshalPropertyMatchers(dbConfig.PropertyMatchers)
	if err != nil {
		return nil, err
	}

	var domainMatchers []domain.PropertyMatcher
	for _, matcher := range matchers {
		domainMatcher, err := domain.NewPropertyMatcher(matcher.Name, matcher.Value)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create property matchers")
		}
		domainMatchers = append(domainMatchers, *domainMatcher)
	}

	return domain.NewIdentifierConfiguration(dbConfig.Id, dbConfig.TargetBusinessObject, domainMatchers, dbConfig.Counter, dbConfig.Prefix, dbConfig.Suffix)
}

func unmarshalToIdentConfigQuery(dbConfig identifierConfiguration) (query.IdentifierConfiguration, error) {
	matchers, err := unmarshalPropertyMatchers(dbConfig.PropertyMatchers)
	if err != nil {
		return query.IdentifierConfiguration{}, err
	}

	var queryMatchers []query.PropertyMatcher
	for _, dbMatcher := range matchers {
		queryMatchers = append(queryMatchers, query.PropertyMatcher{
			Name:  dbMatcher.Name,
			Value: dbMatcher.Value,
		})
	}
	return query.IdentifierConfiguration{
		Id:                   dbConfig.Id,
		TargetBusinessObject: dbConfig.TargetBusinessObject,
		PropertyMatchers:     queryMatchers,
		Counter:              dbConfig.Counter,
		Prefix:               dbConfig.Prefix,
		Suffix:               dbConfig.Suffix,
		CreatedAt:            dbConfig.CreatedAt,
		UpdatedAt:            dbConfig.UpdatedAt,
	}, nil
}

func unmarshalToIdentifier(dbIdentifier identifier) query.Identifier {
	return query.Identifier{
		Id:                        dbIdentifier.Id,
		IdentifierConfigurationId: dbIdentifier.IdentifierConfigurationId,
		Identifier:                dbIdentifier.Identifier,
		CreatedAt:                 dbIdentifier.CreatedAt,
	}
}
