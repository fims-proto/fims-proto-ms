package query

import (
	"context"

	"github.com/google/uuid"
)

type UsersReadModel interface {
	ReadById(ctx context.Context, id uuid.UUID) (User, error)
	ReadByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]User, error)
}

type ReadUsersHandler struct {
	readModel UsersReadModel
}

func NewReadUsersHandler(readModel UsersReadModel) ReadUsersHandler {
	if readModel == nil {
		panic("nil users read model")
	}
	return ReadUsersHandler{readModel: readModel}
}

func (h ReadUsersHandler) HandleReadById(ctx context.Context, id uuid.UUID) (User, error) {
	return h.readModel.ReadById(ctx, id)
}

func (h ReadUsersHandler) HandleReadByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]User, error) {
	return h.readModel.ReadByIds(ctx, ids)
}
