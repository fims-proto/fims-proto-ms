package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain"
)

type CreateIdentifierConfigurationCmd struct {
	Id                   uuid.UUID
	TargetBusinessObject string
	PropertyMatchers     []struct{ Name, Value string }
	Prefix               string
	Suffix               string
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
	var propertyMatchers []domain.PropertyMatcher
	for _, matcher := range cmd.PropertyMatchers {
		propertyMatcher, err := domain.NewPropertyMatcher(matcher.Name, matcher.Value)
		if err != nil {
			return errors.Wrap(err, "failed to handle configuration identifier creation")
		}
		propertyMatchers = append(propertyMatchers, *propertyMatcher)
	}

	configuration, err := domain.NewIdentifierConfiguration(cmd.Id, cmd.TargetBusinessObject, propertyMatchers, 0, cmd.Prefix, cmd.Suffix)
	if err != nil {
		return errors.Wrap(err, "failed to handle configuration identifier creation")
	}

	return h.repo.CreateIdentifierConfiguration(ctx, configuration)
}
