package db

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier_configuration"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/numbering/app/query"
)

type propertyMatcher struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type identifierConfigurationPO struct {
	Id                   uuid.UUID    `gorm:"type:uuid"`
	TargetBusinessObject string       `gorm:"uniqueIndex:identifierConfigs_target_matcher_key"`
	PropertyMatchers     pgtype.JSONB `gorm:"uniqueIndex:identifierConfigs_target_matcher_key"`
	Counter              int
	Prefix               string
	Suffix               string
	CreatedAt            time.Time `gorm:"<-:create"`
	UpdatedAt            time.Time
}

type identifierPO struct {
	Id                        uuid.UUID `gorm:"type:uuid"`
	IdentifierConfigurationId uuid.UUID `gorm:"type:uuid;uniqueIndex:identifiers_configuration_identifier_key"`
	Identifier                string    `gorm:"uniqueIndex:identifiers_configuration_identifier_key"`
	CreatedAt                 time.Time `gorm:"<-:create"`
}

// table names

func (c identifierConfigurationPO) TableName() string {
	return "a_identifier_configurations"
}

func (i identifierPO) TableName() string {
	return "a_identifiers"
}

// mappers

func identifierConfigurationBOToPO(bo identifier_configuration.IdentifierConfiguration) (identifierConfigurationPO, error) {
	matcherPO, err := serializePropertyMatchers(bo.PropertyMatchers())
	if err != nil {
		return identifierConfigurationPO{}, err
	}

	return identifierConfigurationPO{
		Id:                   bo.Id(),
		TargetBusinessObject: bo.TargetBusinessObject(),
		PropertyMatchers:     matcherPO,
		Counter:              bo.Counter(),
		Prefix:               bo.Prefix(),
		Suffix:               bo.Suffix(),
	}, nil
}

func identifierBOToPO(bo identifier.Identifier) identifierPO {
	return identifierPO{
		Id:                        bo.Id(),
		IdentifierConfigurationId: bo.IdentifierConfigurationId(),
		Identifier:                bo.Identifier(),
	}
}

func identifierConfigurationPOToBO(po identifierConfigurationPO) (*identifier_configuration.IdentifierConfiguration, error) {
	matchers, err := deserializePropertyMatchers(po.PropertyMatchers)
	if err != nil {
		return nil, err
	}

	var matcherBOs []identifier_configuration.PropertyMatcher
	for _, matcher := range matchers {
		matcherBO, err := identifier_configuration.NewPropertyMatcher(matcher.Name, matcher.Value)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create property matchers")
		}
		matcherBOs = append(matcherBOs, *matcherBO)
	}

	return identifier_configuration.New(
		po.Id,
		po.TargetBusinessObject,
		matcherBOs,
		po.Counter,
		po.Prefix,
		po.Suffix,
	)
}

func identifierConfigurationPOToDTO(po identifierConfigurationPO) (query.IdentifierConfiguration, error) {
	matchers, err := deserializePropertyMatchers(po.PropertyMatchers)
	if err != nil {
		return query.IdentifierConfiguration{}, err
	}

	var matcherDTOs []query.PropertyMatcher
	for _, matcher := range matchers {
		matcherDTOs = append(matcherDTOs, query.PropertyMatcher{
			Name:  matcher.Name,
			Value: matcher.Value,
		})
	}

	return query.IdentifierConfiguration{
		Id:                   po.Id,
		TargetBusinessObject: po.TargetBusinessObject,
		PropertyMatchers:     matcherDTOs,
		Counter:              po.Counter,
		Prefix:               po.Prefix,
		Suffix:               po.Suffix,
		CreatedAt:            po.CreatedAt,
		UpdatedAt:            po.UpdatedAt,
	}, nil
}

func identifierPOToDTO(po identifierPO) query.Identifier {
	return query.Identifier(po)
}

func serializePropertyMatchers(matcherBOs []identifier_configuration.PropertyMatcher) (pgtype.JSONB, error) {
	var matchers []propertyMatcher
	for _, matcher := range matcherBOs {
		matchers = append(matchers, propertyMatcher{
			Name:  matcher.Name(),
			Value: matcher.Value(),
		})
	}
	var matcherPO pgtype.JSONB
	if err := matcherPO.Set(matchers); err != nil {
		return pgtype.JSONB{}, errors.Wrapf(err, "failed to convert %v to pgtype.JSONB", matcherBOs)
	}

	return matcherPO, nil
}

func deserializePropertyMatchers(po pgtype.JSONB) ([]propertyMatcher, error) {
	var matchers []propertyMatcher
	if err := po.AssignTo(&matchers); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal property matchers")
	}
	return matchers, nil
}
