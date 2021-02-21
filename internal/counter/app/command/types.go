package command

type CounterNextCmd struct {
	UUID string
}

type CounterResetCmd struct {
	UUID string
}

type CounterDeleteCmd struct {
	UUID string
}

type CounterAddCmd struct {
	UUID   string
	Length uint
	Prefix string
	Sufix  string
}
