package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/numbering/app/query"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain"
)

type GenerateNextIdentifierCmd struct {
	IdentifierId         uuid.UUID
	TargetBusinessObject string
	ObjectsToMatch       map[string]string
}

type GenerateNextIdentifierHandler struct {
	repo               domain.Repository
	resolveIdentConfig query.ResolveIdentifierConfigurationReadModel
}

func NewGenerateNextIdentifierHandler(repo domain.Repository, resolveIdentConfig query.ResolveIdentifierConfigurationReadModel) GenerateNextIdentifierHandler {
	if repo == nil {
		panic("nil numbering repo")
	}
	if resolveIdentConfig == nil {
		panic("nil ResolveIdentifierConfigurationReadModel")
	}

	return GenerateNextIdentifierHandler{
		repo:               repo,
		resolveIdentConfig: resolveIdentConfig,
	}
}

func (h GenerateNextIdentifierHandler) Handle(ctx context.Context, cmd GenerateNextIdentifierCmd) error {
	configuration, err := h.resolveIdentConfig.ResolveIdentifierConfiguration(ctx, cmd.TargetBusinessObject, cmd.ObjectsToMatch)
	if err != nil {
		return errors.Wrap(err, "failed to handle generate identifier")
	}

	return h.repo.UpdateIdentifierConfiguration(
		ctx,
		configuration.Id,
		func(config *domain.IdentifierConfiguration) (*domain.IdentifierConfiguration, error) {
			config.IncrementCounter()

			identifier, err := domain.NewIdentifier(cmd.IdentifierId, config.Id(), config.Stringify())
			if err != nil {
				return nil, errors.Wrap(err, "failed to create identifier domain entity")
			}

			err = h.repo.CreateIdentifier(ctx, identifier)
			if err != nil {
				return nil, errors.Wrap(err, "failed to create identifier")
			}

			return config, nil
		},
	)
}
