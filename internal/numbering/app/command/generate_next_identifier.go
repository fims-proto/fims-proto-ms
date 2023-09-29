package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain/identifier_configuration"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/numbering/app/query"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain"
)

type GenerateNextIdentifierCmd struct {
	IdentifierId         uuid.UUID
	TargetBusinessObject string
	ObjectsToMatch       map[string]string
}

type GenerateNextIdentifierHandler struct {
	repo      domain.Repository
	readModel query.NumberingReadModel
}

func NewGenerateNextIdentifierHandler(repo domain.Repository, readModel query.NumberingReadModel) GenerateNextIdentifierHandler {
	if repo == nil {
		panic("nil numbering repo")
	}

	if readModel == nil {
		panic("nil ResolveIdentifierConfigurationReadModel")
	}

	return GenerateNextIdentifierHandler{
		repo:      repo,
		readModel: readModel,
	}
}

func (h GenerateNextIdentifierHandler) Handle(ctx context.Context, cmd GenerateNextIdentifierCmd) error {
	configuration, err := h.readModel.ResolveIdentifierConfiguration(ctx, cmd.TargetBusinessObject, cmd.ObjectsToMatch)
	if err != nil {
		return fmt.Errorf("failed to handle generate identifier: %w", err)
	}

	return h.repo.UpdateIdentifierConfiguration(
		ctx,
		configuration.Id,
		func(config *identifier_configuration.IdentifierConfiguration) (*identifier_configuration.IdentifierConfiguration, error) {
			config.IncrementCounter()

			identifierBO, err := identifier.New(cmd.IdentifierId, config.Id(), config.Stringify())
			if err != nil {
				return nil, fmt.Errorf("failed to create identifier domain entity: %w", err)
			}

			err = h.repo.CreateIdentifier(ctx, identifierBO)
			if err != nil {
				return nil, fmt.Errorf("failed to create identifier: %w", err)
			}

			return config, nil
		},
	)
}
