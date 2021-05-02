package app

import (
	"github/fims-proto/fims-proto-ms/internal/counter/app/command"
	"github/fims-proto/fims-proto-ms/internal/counter/app/query"
	"github/fims-proto/fims-proto-ms/internal/counter/domain"
)

type Commands struct {
	NextCounter   command.CounterNextHandler
	DeleteCounter command.CounterDeleteHandler
	ResetCounter  command.CounterResetHandler
	CreateCounter command.CounterCreateHandler
	LoadCounters  command.CounterDataloadHandler
}

type Queries struct {
	ReadCounters query.ReadCountersHandler
}

type Application struct {
	Commands Commands
	Queries  Queries
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(readModel query.CountersReadModel, repo domain.Repository) {
	a.Queries = Queries{
		ReadCounters: query.NewReadCountersHandler(readModel),
	}
	a.Commands = Commands{
		NextCounter:   command.NewCounterNextHandler(repo),
		DeleteCounter: command.NewCounterDeleteHandler(repo),
		ResetCounter:  command.NewCounterResetHandler(repo),
		CreateCounter: command.NewCounterCreateHandler(repo),
		LoadCounters:  command.NewCounterDataloadHandler(repo),
	}
}
