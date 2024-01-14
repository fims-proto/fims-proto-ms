package template

import (
	reportRule "github/fims-proto/fims-proto-ms/internal/report/domain/template/rule"

	"github.com/google/uuid"
)

type Formula struct {
	accountId        uuid.UUID
	itemId           uuid.UUID
	isAccountFormula bool
	sumFactor        int
	rule             reportRule.Rule
}
