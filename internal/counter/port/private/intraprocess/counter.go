package intraprocess

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/counter/app"
	"github/fims-proto/fims-proto-ms/internal/counter/app/command"
)

type CounterInterface struct {
	app app.Application
}

func NewHandler(app app.Application) CounterInterface {
	return CounterInterface{app: app}
}

func (h CounterInterface) Next(ctx context.Context, UUID string) (string, error) {
	return h.app.Commands.NextCounter.Handle(ctx, command.CounterNextCmd{UUID: UUID})
}

func (h CounterInterface) Reset(ctx context.Context, UUID string) error {
	return h.app.Commands.ResetCounter.Handle(ctx, command.CounterResetCmd{UUID: UUID})
}

func (h CounterInterface) Delete(ctx context.Context, UUID string) error {
	return h.app.Commands.DeleteCounter.Handle(ctx, command.CounterDeleteCmd{UUID: UUID})
}

func (h CounterInterface) Add(ctx context.Context, UUID string, len uint, prefix string, sufix string) error {
	return h.app.Commands.AddCounter.Handle(
		ctx,
		command.CounterAddCmd{
			UUID:   UUID,
			Length: len,
			Prefix: prefix,
			Sufix:  sufix,
		})
}
