package intraprocess

import "github/fims-proto/fims-proto-ms/internal/counter/app/command"

func (req CreateCounterRequest) mapToCommand() command.CounterCreateCmd {
	return command.CounterCreateCmd{
		Prefix:          req.Prefix,
		Suffix:          req.Suffix,
		BusinessObjects: req.BusinessObjects,
	}
}
