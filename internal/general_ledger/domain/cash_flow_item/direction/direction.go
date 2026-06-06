package direction

import "fmt"

type Direction struct {
	slug string
}

func (d Direction) String() string {
	return d.slug
}

var (
	unknown = Direction{""}
	Inflow  = Direction{"INFLOW"}
	Outflow = Direction{"OUTFLOW"}
)

func FromString(s string) (Direction, error) {
	switch s {
	case Inflow.slug:
		return Inflow, nil
	case Outflow.slug:
		return Outflow, nil
	}
	return unknown, fmt.Errorf("unknown cash flow direction: %s", s)
}
