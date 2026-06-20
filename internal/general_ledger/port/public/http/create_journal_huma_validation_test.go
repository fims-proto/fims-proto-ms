package http

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestCreateJournalHumaValidationAcceptsClientPayloadShape(t *testing.T) {
	_, api := humatest.New(t, huma.DefaultConfig("Test API", "1.0.0"))
	huma.Register(api, huma.Operation{
		Method: "POST",
		Path:   "/api/v1/sob/{sobId}/journals",
	}, func(ctx context.Context, input *CreateJournalInput) (*createJournalValidationOutput, error) {
		return &createJournalValidationOutput{Status: http.StatusNoContent}, nil
	})

	resp := api.Post(
		"/api/v1/sob/11111111-1111-1111-1111-111111111111/journals",
		"Content-Type: application/json",
		strings.NewReader(`{
			"headerText": "test journal",
			"attachmentQuantity": 0,
			"creator": "11111111-1111-1111-1111-111111111111",
			"transactionDate": "2026-06-19",
			"journalLines": [
				{
					"amount": 111,
					"rawAccountNumber": "001001",
					"text": "111"
				},
				{
					"amount": -111,
					"cashFlowItemId": "7caf2919-4756-4ded-bcff-c870fe6e827b",
					"rawAccountNumber": "001101000001",
					"text": "111"
				}
			]
		}`),
	)

	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected Huma to accept payload and reach handler, got status %d: %s", resp.Code, resp.Body.String())
	}
}

type createJournalValidationOutput struct {
	Status int
}
