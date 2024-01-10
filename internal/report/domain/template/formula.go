package report

import (
	reportRule "github/fims-proto/fims-proto-ms/internal/report/domain/template/rule"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Formula struct {
	accountId          uuid.UUID
	sumFactor          int
	rule               reportRule.Rule
	values             []Cell[decimal.Decimal]
	displayAsBreakdown bool
}
