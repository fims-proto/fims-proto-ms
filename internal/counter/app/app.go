package app

import "github/fims-proto/fims-proto-ms/internal/counter/app/command"

type Commands struct {
	NextCounter   command.CounterNextHandler
	DeleteCounter command.CounterDeleteHandler
	ResetCounter  command.CounterResetHandler
	AddCounter    command.CounterAddHandler
}

type Application struct {
	Commands Commands
	// Queries currently empty
}
