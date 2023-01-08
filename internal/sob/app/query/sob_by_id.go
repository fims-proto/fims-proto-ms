package query

import (
	"context"

	"github.com/google/uuid"
)

type SobByIdHandler struct {
	readModel SobReadModel
}

func NewSobByIdHandler(readModel SobReadModel) SobByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return SobByIdHandler{readModel: readModel}
}

func (h SobByIdHandler) Handle(ctx context.Context, sobId uuid.UUID) (Sob, error) {
	return h.readModel.SobById(ctx, sobId)
}
