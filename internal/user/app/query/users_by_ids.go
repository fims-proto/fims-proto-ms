package query

import (
	"context"

	"github.com/google/uuid"
)

type UsersByIdsHandler struct {
	readModel UserReadModel
}

func NewUsersByIdsHandler(readModel UserReadModel) UsersByIdsHandler {
	if readModel == nil {
		panic("nil users read model")
	}

	return UsersByIdsHandler{readModel: readModel}
}

func (h UsersByIdsHandler) Handle(ctx context.Context, userIds []uuid.UUID) ([]User, error) {
	return h.readModel.UsersByIds(ctx, userIds)
}
