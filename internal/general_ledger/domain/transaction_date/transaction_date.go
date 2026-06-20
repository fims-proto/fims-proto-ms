package transaction_date

import (
	"encoding/json"
	"fmt"
)

// TransactionDate represents a business date without timezone semantics.
// It represents a calendar date (year, month, day) that has the same meaning
// regardless of timezone.
type TransactionDate struct {
	Year  int
	Month int
	Day   int
}

func (td TransactionDate) format() string {
	return fmt.Sprintf("%04d-%02d-%02d", td.Year, td.Month, td.Day)
}

// IsZero returns true if this is the zero value of TransactionDate.
func (td TransactionDate) IsZero() bool {
	return td.Year == 0 || td.Month == 0 || td.Day == 0
}

// MarshalJSON implements json.Marshaler to serialize as ISO 8601 string.
// Output format: "2026-02-01"
func (td TransactionDate) MarshalJSON() ([]byte, error) {
	if td.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(td.format())
}
