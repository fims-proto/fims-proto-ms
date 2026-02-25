package transaction_date

import (
	"encoding/json"
	"fmt"
	"time"
)

// TransactionDate represents a business date without timezone semantics.
// It represents a calendar date (year, month, day) that has the same meaning
// regardless of timezone.
type TransactionDate struct {
	Year  int
	Month int
	Day   int
}

// NewTransactionDate creates a new TransactionDate with validation.
// It validates that the date is valid (e.g., not Feb 30, not month 13)
// by attempting to parse it as a time.Time.
func NewTransactionDate(year, month, day int) (TransactionDate, error) {
	// Validate by creating a time.Time and checking if it matches input
	dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	parsedTime, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return TransactionDate{}, fmt.Errorf("invalid date: %w", err)
	}

	// Verify the parsed date matches the input (catches invalid dates like Feb 30)
	if parsedTime.Year() != year || int(parsedTime.Month()) != month || parsedTime.Day() != day {
		return TransactionDate{}, fmt.Errorf("invalid date: %04d-%02d-%02d", year, month, day)
	}

	return TransactionDate{
		Year:  year,
		Month: month,
		Day:   day,
	}, nil
}

// String returns the date in ISO 8601 format (YYYY-MM-DD).
func (td TransactionDate) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", td.Year, td.Month, td.Day)
}

// IsZero returns true if this is the zero value of TransactionDate.
func (td TransactionDate) IsZero() bool {
	return td.Year == 0 || td.Month == 0 || td.Day == 0
}

// Equal returns true if this date equals the other date.
func (td TransactionDate) Equal(other TransactionDate) bool {
	return td.Year == other.Year && td.Month == other.Month && td.Day == other.Day
}

// MarshalJSON implements json.Marshaler to serialize as ISO 8601 string.
// Output format: "2026-02-01"
func (td TransactionDate) MarshalJSON() ([]byte, error) {
	if td.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(td.String())
}

// UnmarshalJSON implements json.Unmarshaler to deserialize from ISO 8601 string.
// Expected format: "2026-02-01"
func (td *TransactionDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("transaction date must be a string: %w", err)
	}

	if s == "" || s == "null" {
		*td = TransactionDate{}
		return nil
	}

	// Parse using time.Parse for automatic validation
	parsedTime, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	td.Year = parsedTime.Year()
	td.Month = int(parsedTime.Month())
	td.Day = parsedTime.Day()

	return nil
}
