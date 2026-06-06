package http

import (
	"encoding/json"
	"testing"

	"github/fims-proto/fims-proto-ms/internal/report/app/query"

	"github.com/stretchr/testify/require"
)

func TestRowDTOToVOEmitsDisplaySumFactor(t *testing.T) {
	vo := rowDTOToVO(query.Row{
		RowCode:          "ROW",
		Text:             "Row",
		ShowLineNo:       true,
		SumFactor:        -1,
		DisplaySumFactor: true,
		Expression:       query.Expression{Kind: "none"},
	})

	require.True(t, vo.DisplaySumFactor)
}

func TestRowDTOToVOEmitsFalseDisplaySumFactor(t *testing.T) {
	vo := rowDTOToVO(query.Row{
		RowCode:    "ROW",
		Text:       "Row",
		ShowLineNo: true,
		SumFactor:  1,
		Expression: query.Expression{Kind: "none"},
	})

	raw, err := json.Marshal(vo)

	require.NoError(t, err)
	require.Contains(t, string(raw), `"displaySumFactor":false`)
}
