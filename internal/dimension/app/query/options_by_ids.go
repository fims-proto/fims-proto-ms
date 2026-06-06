package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type OptionsByIdsHandler struct {
	readModel DimensionReadModel
}

func NewOptionsByIdsHandler(readModel DimensionReadModel) OptionsByIdsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return OptionsByIdsHandler{readModel: readModel}
}

func (h OptionsByIdsHandler) Handle(ctx context.Context, optionIds []uuid.UUID) ([]DimensionOption, error) {
	options, err := h.readModel.OptionsByIds(ctx, optionIds)
	if err != nil {
		return nil, fmt.Errorf("failed to read options by ids: %w", err)
	}

	return options, nil
}
