package request

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/shopspring/decimal"
)

type DateRequest struct {
	Year  int
	Month int
	Day   int
}

func (r *DateRequest) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("date must be a string: %w", err)
	}

	if s == "" || s == "null" {
		*r = DateRequest{}
		return nil
	}

	parsed, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	r.Year = parsed.Year()
	r.Month = int(parsed.Month())
	r.Day = parsed.Day()
	return nil
}

func (*DateRequest) Schema(_ huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type:    huma.TypeString,
		Format:  "date",
		Pattern: `^\d{4}-\d{2}-\d{2}$`,
	}
}

type DecimalRequest decimal.Decimal

func (r *DecimalRequest) UnmarshalJSON(data []byte) error {
	var amount decimal.Decimal
	if err := amount.UnmarshalJSON(data); err != nil {
		return err
	}
	*r = DecimalRequest(amount)
	return nil
}

func (*DecimalRequest) Schema(_ huma.Registry) *huma.Schema {
	return &huma.Schema{
		AnyOf: []*huma.Schema{
			{Type: huma.TypeNumber},
			{Type: huma.TypeString},
		},
	}
}
