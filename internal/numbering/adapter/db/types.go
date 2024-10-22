package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/numbering/app/query"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier_configuration"
)

type propertyMatcher struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type identifierConfigurationPO struct {
	Id                   uuid.UUID         `gorm:"type:uuid"`
	TargetBusinessObject string            `gorm:"uniqueIndex:UQ_IdentifierConfigurations_TargetBusinessObject_PropertyMatchers"`
	PropertyMatchers     []propertyMatcher `gorm:"type:jsonb;serializer:json;uniqueIndex:UQ_IdentifierConfigurations_TargetBusinessObject_PropertyMatchers"`
	Counter              int
	Prefix               string
	Suffix               string
	CreatedAt            time.Time `gorm:"<-:create"`
	UpdatedAt            time.Time
}

type identifierPO struct {
	Id                        uuid.UUID `gorm:"type:uuid"`
	IdentifierConfigurationId uuid.UUID `gorm:"type:uuid;uniqueIndex:UQ_Identifiers_IdentifierConfigurationId_Identifier"`
	Identifier                string    `gorm:"uniqueIndex:UQ_Identifiers_IdentifierConfigurationId_Identifier"`
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

func propertyMatcherBOToPO(bos []identifier_configuration.PropertyMatcher) []propertyMatcher {
	var result []propertyMatcher
	for _, bo := range bos {
		result = append(result, propertyMatcher{
			Name:  bo.Name(),
			Value: bo.Value(),
		})
	}

	return result
}

func identifierConfigurationBOToPO(bo identifier_configuration.IdentifierConfiguration) (identifierConfigurationPO, error) {
	return identifierConfigurationPO{
		Id:                   bo.Id(),
		TargetBusinessObject: bo.TargetBusinessObject(),
		PropertyMatchers:     propertyMatcherBOToPO(bo.PropertyMatchers()),
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
	var matcherBOs []identifier_configuration.PropertyMatcher
	for _, matcher := range po.PropertyMatchers {
		matcherBO, err := identifier_configuration.NewPropertyMatcher(matcher.Name, matcher.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to create property matchers: %w", err)
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
	var matcherDTOs []query.PropertyMatcher
	for _, matcher := range po.PropertyMatchers {
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
