package http

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUpdateReportRequestMapToCommandIgnoresColumns(t *testing.T) {
	raw := []byte(`{
		"title": "Updated Report",
		"columns": [
			{"id": "not-a-uuid", "label": "", "valueType": "unsupported"}
		],
		"rows": []
	}`)

	var req UpdateReportRequest
	require.NoError(t, json.Unmarshal(raw, &req))

	cmd, err := req.mapToCommand(uuid.New(), uuid.New())

	require.NoError(t, err)
	require.NotNil(t, cmd.Title)
	require.Equal(t, "Updated Report", *cmd.Title)
	require.Empty(t, cmd.Rows)
}

func TestUpdateReportRequestMapToCommandMapsNestedDisplaySumFactor(t *testing.T) {
	raw := []byte(`{
		"rows": [
			{
				"rowCode": "PARENT",
				"text": "Parent",
				"showLineNo": true,
				"sumFactor": 1,
				"expression": {"kind": "none"},
				"rows": [
					{
						"rowCode": "CHILD",
						"text": "Child",
						"showLineNo": true,
						"sumFactor": -1,
						"displaySumFactor": true,
						"expression": {"kind": "none"}
					}
				]
			}
		]
	}`)

	var req UpdateReportRequest
	require.NoError(t, json.Unmarshal(raw, &req))

	cmd, err := req.mapToCommand(uuid.New(), uuid.New())

	require.NoError(t, err)
	require.False(t, cmd.Rows[0].DisplaySumFactor)
	require.True(t, cmd.Rows[0].Rows[0].DisplaySumFactor)
}
