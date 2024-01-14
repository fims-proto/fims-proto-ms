package report

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"
	template "github/fims-proto/fims-proto-ms/internal/report/domain/template" // Import the package that contains the Template struct

	"github.com/google/uuid"
)

type Report struct {
	id                uuid.UUID
	periodId          uuid.UUID
	refTemplateId     uuid.UUID // different from the templateId inside the inner template
	version           int
	template.Template // an inner template deep-copied from the refTemplate
}

func New(
	id uuid.UUID,
	periodId uuid.UUID,
	version int,
	refTemplateId uuid.UUID,
	inner *template.Template,
) (*Report, error) {
	if id == uuid.Nil {
		return nil, errors.NewSlugError("voucher-emptyId")
	}
	if periodId == uuid.Nil {
		return nil, errors.NewSlugError("emptyPeriodId")
	}
	if refTemplateId == uuid.Nil {
		return nil, errors.NewSlugError("emptyTemplateId")
	}
	return &Report{
		id:            id,
		periodId:      periodId,
		refTemplateId: refTemplateId,
		version:       version,
		Template:      *inner,
	}, nil
}
