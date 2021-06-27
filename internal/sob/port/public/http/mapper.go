package http

import (
	"github/fims-proto/fims-proto-ms/internal/sob/app/command"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
)

func mapFromSobQuery(q query.Sob) SobResponse {
	return SobResponse{
		Id:          q.Id,
		Name:        q.Name,
		Description: q.Description,
	}
}

func (r CreateSobRequest) mapToCommand() command.CreateSobCmd {
	return command.CreateSobCmd{
		Id:          r.Id,
		Name:        r.Name,
		Description: r.Description,
	}
}
