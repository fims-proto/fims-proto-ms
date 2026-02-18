package user

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/user/app/query"
	userPort "github/fims-proto/fims-proto-ms/internal/user/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	userInterface userPort.UserInterface
}

func NewIntraProcessAdapter(userInterface userPort.UserInterface) IntraProcessAdapter {
	return IntraProcessAdapter{userInterface: userInterface}
}

func (i IntraProcessAdapter) ReadUsersByIds(ctx context.Context, userIds []uuid.UUID) (map[uuid.UUID]query.User, error) {
	users, err := i.userInterface.ReadUsersByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	usersMap := make(map[uuid.UUID]query.User)
	for _, user := range users {
		usersMap[user.Id] = user
	}

	return usersMap, nil
}
