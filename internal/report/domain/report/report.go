package report

import (
	template "github/fims-proto/fims-proto-ms/internal/report/domain/template" // Import the package that contains the Template struct

	"github.com/google/uuid"
)

type Report struct {
	id                uuid.UUID
	periodId          uuid.UUID
	version           int
	template.Template // Use the imported package to reference the Template struct
}
