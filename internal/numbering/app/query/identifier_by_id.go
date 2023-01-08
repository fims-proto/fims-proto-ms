package query

import (
	"context"

	"github.com/google/uuid"
)

type IdentifierByIdHandler struct {
	readModel NumberingReadModel
}

func NewIdentifierByIdHandler(readModel NumberingReadModel) IdentifierByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return IdentifierByIdHandler{readModel: readModel}
}

func (h IdentifierByIdHandler) Handle(ctx context.Context, id uuid.UUID) (Identifier, error) {
	return h.readModel.IdentifierById(ctx, id)
}
