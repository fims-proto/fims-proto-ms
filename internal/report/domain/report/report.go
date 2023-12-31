package report

import "github.com/google/uuid"

type Report struct {
	// id is for report id. Differentiating from templateId in Template struct.
	id       uuid.UUID
	periodId uuid.UUID
	version  int
	Template
}
