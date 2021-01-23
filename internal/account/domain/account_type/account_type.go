package accounttype

import "github.com/pkg/errors"

// enum
type Type struct {
	name string
}

var (
	Assets        = Type{"assets"}
	Liabilities   = Type{"liabilities"}
	Common        = Type{"common"}
	Equity        = Type{"equity"}
	Cost          = Type{"cost"}
	ProfitAndLoss = Type{"profit_and_loss"}
)

var availableTypes = []Type{
	Assets,
	Liabilities,
	Common,
	Equity,
	Cost,
	ProfitAndLoss,
}

func NewAccountTypeFromString(name string) (Type, error) {
	for _, accType := range availableTypes {
		if accType.String() == name {
			return accType, nil
		}
	}

	return Type{}, errors.Errorf("unknown Account Type: '%s'", name)
}

func (accType Type) String() string {
	return accType.name
}
