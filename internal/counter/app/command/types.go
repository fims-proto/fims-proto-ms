package command

import "github.com/google/uuid"

type CounterNextCmd struct {
	CounterUUID uuid.UUID
}

type CounterResetCmd struct {
	CounterUUID uuid.UUID
}

type CounterDeleteCmd struct {
	CounterUUID uuid.UUID
}

type CounterCreateCmd struct {
	Prefix string
	Sufix  string
}
