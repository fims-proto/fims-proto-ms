package intraprocess

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/counter/app"
	"github/fims-proto/fims-proto-ms/internal/counter/app/command"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CounterInterface struct {
	app app.Application
}

func NewCounterInterface(app app.Application) CounterInterface {
	return CounterInterface{app: app}
}

func (i CounterInterface) Create(ctx context.Context, req CreateCounterRequest) (uuid.UUID, error) {
	return i.app.Commands.CreateCounter.Handle(ctx, req.mapToCommand())
}

func (i CounterInterface) Next(ctx context.Context, businessObject string) (string, error) {
	return i.queryThenProceed(ctx, businessObject, func(ctx context.Context, counterUUID uuid.UUID) (string, error) {
		return i.app.Commands.NextCounter.Handle(ctx, command.CounterNextCmd{CounterUUID: counterUUID})
	})
}

func (i CounterInterface) Reset(ctx context.Context, businessObject string) error {
	_, err := i.queryThenProceed(ctx, businessObject, func(ctx context.Context, counterUUID uuid.UUID) (string, error) {
		return "DUMMY", i.app.Commands.ResetCounter.Handle(ctx, command.CounterResetCmd{CounterUUID: counterUUID})
	})
	return err
}

func (i CounterInterface) Delete(ctx context.Context, businessObject string) error {
	_, err := i.queryThenProceed(ctx, businessObject, func(ctx context.Context, counterUUID uuid.UUID) (string, error) {
		return "DUMMY", i.app.Commands.DeleteCounter.Handle(ctx, command.CounterDeleteCmd{CounterUUID: counterUUID})
	})
	return err
}

func (i CounterInterface) queryThenProceed(
	ctx context.Context,
	businessObject string,
	proceed func(ctx context.Context, counterUUID uuid.UUID) (string, error),
) (string, error) {
	counter, err := i.app.Queries.ReadCounters.HandleByBusinessObject(ctx, businessObject)
	if err != nil {
		return "", errors.Wrapf(err, "read counter failed with business object %s", businessObject)
	}
	return proceed(ctx, counter.CounterUUID)
}
