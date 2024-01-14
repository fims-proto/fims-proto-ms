package template

import (
	"github/fims-proto/fims-proto-ms/internal/common/errors"
	templateClass "github/fims-proto/fims-proto-ms/internal/report/domain/template/class"

	"github.com/google/uuid"
)

type Template struct {
	// Template will be embedded directly into Report, so here name it as templateId. There will be another id field in Report.
	templateId uuid.UUID
	sobId      uuid.UUID
	class      templateClass.Class
	title      string
	tables     []*Table
}

func New(
	id uuid.UUID,
	sobId uuid.UUID,
	class string,
	title string,
	tables []*Table,
) (*Template, error) {
	if id == uuid.Nil {
		return nil, errors.NewSlugError("voucher-emptyId")
	}

	if sobId == uuid.Nil {
		return nil, errors.NewSlugError("emptySobId")
	}

	if title == "" {
		return nil, errors.NewSlugError("template-emptyTitle")
	}
	cl, err := templateClass.FromString(class)
	if err != nil {
		return nil, err
	}

	return &Template{
		templateId: id,
		sobId:      sobId,
		class:      cl,
		tables:     tables,
	}, nil
}
