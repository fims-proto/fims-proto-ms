package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier_configuration"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain"
)

type CreateIdentifierConfigurationCmd struct {
	IdentifierConfigurationId uuid.UUID
	TargetBusinessObject      string
	PropertyMatchers          []struct{ Name, Value string }
	Prefix                    string
	Suffix                    string
}

type CreateIdentifierConfigurationHandler struct {
	repo domain.Repository
}

func NewCreateIdentifierConfigurationHandler(repo domain.Repository) CreateIdentifierConfigurationHandler {
	if repo == nil {
		panic("nil numbering repo")
	}

	return CreateIdentifierConfigurationHandler{repo: repo}
}

func (h CreateIdentifierConfigurationHandler) Handle(ctx context.Context, cmd CreateIdentifierConfigurationCmd) error {
	var propertyMatchers []identifier_configuration.PropertyMatcher
	for _, matcher := range cmd.PropertyMatchers {
		propertyMatcher, err := identifier_configuration.NewPropertyMatcher(matcher.Name, matcher.Value)
		if err != nil {
			return fmt.Errorf("failed to handle configuration identifier creation: %w", err)
		}
		propertyMatchers = append(propertyMatchers, *propertyMatcher)
	}

	configuration, err := identifier_configuration.New(cmd.IdentifierConfigurationId, cmd.TargetBusinessObject, propertyMatchers, 0, cmd.Prefix, cmd.Suffix)
	if err != nil {
		return fmt.Errorf("failed to handle configuration identifier creation: %w", err)
	}

	return h.repo.CreateIdentifierConfiguration(ctx, configuration)
}
