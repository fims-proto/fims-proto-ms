package period

import "github/fims-proto/fims-proto-ms/internal/common/errors"

func (p *Period) Start() error {
	if p.isClosed {
		return errors.NewInvalidInputError(errors.SlugPeriodCloseClosed)
	}

	p.isCurrent = true
	return nil
}
