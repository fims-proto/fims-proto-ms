package intraprocess

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/counter/app"
	"github/fims-proto/fims-proto-ms/internal/counter/app/command"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CounterInterface struct {
	app *app.Application
}

func NewCounterInterface(app *app.Application) CounterInterface {
	return CounterInterface{app: app}
}

func (i CounterInterface) InitializeCounters(ctx context.Context, sobId uuid.UUID) error {
	return i.app.Commands.LoadCounters.Handle(ctx, sobId)
}

func (i CounterInterface) Create(ctx context.Context, req CreateCounterRequest) error {
	return i.app.Commands.CreateCounter.Handle(ctx, req.mapToCommand())
}

func (i CounterInterface) Next(ctx context.Context, sep string, businessObjects []string) (string, error) {
	return i.queryThenProceed(
		ctx,
		sep,
		businessObjects,
		func(ctx context.Context, counterUUID uuid.UUID) (string, error) {
			return i.app.Commands.NextCounter.Handle(ctx, command.CounterNextCmd{CounterUUID: counterUUID})
		},
	)
}

func (i CounterInterface) Reset(ctx context.Context, sep string, businessObjects []string) error {
	_, err := i.queryThenProceed(
		ctx,
		sep,
		businessObjects,
		func(ctx context.Context, counterUUID uuid.UUID) (string, error) {
			return "DUMMY", i.app.Commands.ResetCounter.Handle(ctx, command.CounterResetCmd{CounterUUID: counterUUID})
		},
	)
	return err
}

func (i CounterInterface) Delete(ctx context.Context, sep string, businessObjects []string) error {
	_, err := i.queryThenProceed(
		ctx,
		sep,
		businessObjects,
		func(ctx context.Context, counterUUID uuid.UUID) (string, error) {
			return "DUMMY", i.app.Commands.DeleteCounter.Handle(ctx, command.CounterDeleteCmd{CounterUUID: counterUUID})
		},
	)
	return err
}

func (i CounterInterface) queryThenProceed(
	ctx context.Context,
	sep string,
	businessObjects []string,
	proceed func(ctx context.Context, counterUUID uuid.UUID) (string, error),
) (string, error) {
	counter, err := i.app.Queries.ReadCounters.HandleByBusinessObject(ctx, sep, businessObjects)
	if err != nil {
		return "", errors.Wrapf(err, "read counter failed with business object: separator %s, objects %s", sep, businessObjects)
	}
	return proceed(ctx, counter.Id)
}
