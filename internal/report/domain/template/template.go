package report

import (
	reportClass "github/fims-proto/fims-proto-ms/internal/report/domain/template/class"

	"github.com/google/uuid"
)

type Template struct {
	// Template will be embedded directly into Report, so here name it as templateId. There will be another id field in Report.
	templateId uuid.UUID
	sobId      uuid.UUID
	class      reportClass.Class
	title      string
	tables     []Table
}
