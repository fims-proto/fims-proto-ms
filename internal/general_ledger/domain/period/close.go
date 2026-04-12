package period

import "github/fims-proto/fims-proto-ms/internal/common/errors"

func (p *Period) Close() error {
	if p.isClosed {
		return errors.NewInvalidInputError(errors.SlugPeriodCloseClosed)
	}

	if !p.isCurrent {
		return errors.NewInvalidInputError(errors.SlugPeriodCloseIsNotCurrent)
	}

	p.isClosed = true
	p.isCurrent = false
	return nil
}
