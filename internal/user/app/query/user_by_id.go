package query

import (
	"context"

	"github.com/google/uuid"
)

type UserByIdHandler struct {
	readModel UserReadModel
}

func NewUserByIdHandler(readModel UserReadModel) UserByIdHandler {
	if readModel == nil {
		panic("nil users read model")
	}

	return UserByIdHandler{readModel: readModel}
}

func (h UserByIdHandler) Handle(ctx context.Context, userId uuid.UUID) (User, error) {
	return h.readModel.UserById(ctx, userId)
}
