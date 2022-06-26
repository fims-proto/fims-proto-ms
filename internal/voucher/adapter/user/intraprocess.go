package user

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/user/app/query"
	userPort "github/fims-proto/fims-proto-ms/internal/user/port/private/intraprocess"
)

type IntraProcessAdapter struct {
	userInterface userPort.UserInterface
}

func NewIntraProcessAdapter(userInterface userPort.UserInterface) IntraProcessAdapter {
	return IntraProcessAdapter{userInterface: userInterface}
}

func (i IntraProcessAdapter) ReadUserByIds(ctx context.Context, userIds []uuid.UUID) (map[uuid.UUID]query.User, error) {
	return i.userInterface.ReadUserByIds(ctx, userIds)
}
