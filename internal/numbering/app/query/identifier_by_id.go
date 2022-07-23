package query

import (
	"context"

	"github.com/google/uuid"
)

type IdentifierByIdReadModel interface {
	IdentifierById(ctx context.Context, id uuid.UUID) (Identifier, error)
}

type IdentifierByIdHandler struct {
	readModel IdentifierByIdReadModel
}

func NewIdentifierByIdHandler(readModel IdentifierByIdReadModel) IdentifierByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return IdentifierByIdHandler{readModel: readModel}
}

func (h IdentifierByIdHandler) Handle(ctx context.Context, id uuid.UUID) (Identifier, error) {
	return h.readModel.IdentifierById(ctx, id)
}
