package period

import "github/fims-proto/fims-proto-ms/internal/common/errors"

func (p *Period) Close() error {
	if p.isClosed {
		return errors.NewSlugError("period-close-repeatClose")
	}

	if !p.isCurrent {
		return errors.NewSlugError("period-close-isNotCurrent")
	}

	p.isClosed = true
	return nil
}
