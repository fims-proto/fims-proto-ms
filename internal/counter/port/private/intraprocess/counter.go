package intraprocess

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/counter/app"
	"github/fims-proto/fims-proto-ms/internal/counter/app/command"

	"github.com/google/uuid"
)

type CounterInterface struct {
	app app.Application
}

func NewHandler(app app.Application) CounterInterface {
	return CounterInterface{app: app}
}

func (h CounterInterface) Next(ctx context.Context, counterUUID uuid.UUID) (string, error) {
	return h.app.Commands.NextCounter.Handle(ctx, command.CounterNextCmd{CounterUUID: counterUUID})
}

func (h CounterInterface) Reset(ctx context.Context, counterUUID uuid.UUID) error {
	return h.app.Commands.ResetCounter.Handle(ctx, command.CounterResetCmd{CounterUUID: counterUUID})
}

func (h CounterInterface) Delete(ctx context.Context, counterUUID uuid.UUID) error {
	return h.app.Commands.DeleteCounter.Handle(ctx, command.CounterDeleteCmd{CounterUUID: counterUUID})
}

func (h CounterInterface) Create(ctx context.Context, prefix string, sufix string) (uuid.UUID, error) {
	return h.app.Commands.CreateCounter.Handle(
		ctx,
		command.CounterCreateCmd{
			Prefix: prefix,
			Sufix:  sufix,
		})
}
