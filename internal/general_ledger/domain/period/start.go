package period

import "github/fims-proto/fims-proto-ms/internal/common/errors"

func (p *Period) Start() error {
	if p.isClosed {
		return errors.NewSlugError("period-close-closed")
	}

	p.isCurrent = true
	return nil
}
