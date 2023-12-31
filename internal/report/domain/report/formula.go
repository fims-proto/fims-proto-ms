package report

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	reportRule "github/fims-proto/fims-proto-ms/internal/report/domain/report/rule"
)

type Formula struct {
	accountId uuid.UUID
	sumFactor int
	rule      reportRule.Rule
	values    []Cell[decimal.Decimal]
}
