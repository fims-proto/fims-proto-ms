package data_source

import "fmt"

type DataSource struct {
	slug string
}

func (d DataSource) String() string {
	return d.slug
}

var (
	Unknown  = DataSource{""}
	Empty    = DataSource{"empty"}
	Formulas = DataSource{"formulas"}
	Sum      = DataSource{"sum"}
)

var stringToDataSource = map[string]DataSource{
	"empty":    Empty,
	"formulas": Formulas,
	"sum":      Sum,
}

func FromString(s string) (DataSource, error) {
	dataSource, ok := stringToDataSource[s]
	if ok {
		return dataSource, nil
	}
	return Unknown, fmt.Errorf("unknown data source %s", s)
}
