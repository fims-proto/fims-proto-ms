package transaction_date

import (
	"encoding/json"
	"testing"
)

func TestNewTransactionDate(t *testing.T) {
	tests := []struct {
		name      string
		year      int
		month     int
		day       int
		wantError bool
	}{
		{
			name:      "valid date",
			year:      2026,
			month:     2,
			day:       1,
			wantError: false,
		},
		{
			name:      "valid leap year date",
			year:      2024,
			month:     2,
			day:       29,
			wantError: false,
		},
		{
			name:      "invalid leap year date",
			year:      2026,
			month:     2,
			day:       29,
			wantError: true,
		},
		{
			name:      "invalid month",
			year:      2026,
			month:     13,
			day:       1,
			wantError: true,
		},
		{
			name:      "invalid day",
			year:      2026,
			month:     2,
			day:       30,
			wantError: true,
		},
		{
			name:      "zero month",
			year:      2026,
			month:     0,
			day:       1,
			wantError: true,
		},
		{
			name:      "negative day",
			year:      2026,
			month:     1,
			day:       -1,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTransactionDate(tt.year, tt.month, tt.day)
			if (err != nil) != tt.wantError {
				t.Errorf("NewTransactionDate() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError {
				if got.Year != tt.year || got.Month != tt.month || got.Day != tt.day {
					t.Errorf("NewTransactionDate() = %v, want year=%d month=%d day=%d",
						got, tt.year, tt.month, tt.day)
				}
			}
		})
	}
}

func TestTransactionDate_String(t *testing.T) {
	tests := []struct {
		name string
		td   TransactionDate
		want string
	}{
		{
			name: "normal date",
			td:   TransactionDate{Year: 2026, Month: 2, Day: 1},
			want: "2026-02-01",
		},
		{
			name: "single digit month and day",
			td:   TransactionDate{Year: 2026, Month: 1, Day: 5},
			want: "2026-01-05",
		},
		{
			name: "zero value",
			td:   TransactionDate{},
			want: "0000-00-00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.td.String(); got != tt.want {
				t.Errorf("TransactionDate.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionDate_IsZero(t *testing.T) {
	tests := []struct {
		name string
		td   TransactionDate
		want bool
	}{
		{
			name: "zero value",
			td:   TransactionDate{},
			want: true,
		},
		{
			name: "non-zero value",
			td:   TransactionDate{Year: 2026, Month: 2, Day: 1},
			want: false,
		},
		{
			name: "partially zero",
			td:   TransactionDate{Year: 2026, Month: 0, Day: 0},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.td.IsZero(); got != tt.want {
				t.Errorf("TransactionDate.IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionDate_Equal(t *testing.T) {
	tests := []struct {
		name  string
		td    TransactionDate
		other TransactionDate
		want  bool
	}{
		{
			name:  "equal dates",
			td:    TransactionDate{Year: 2026, Month: 2, Day: 1},
			other: TransactionDate{Year: 2026, Month: 2, Day: 1},
			want:  true,
		},
		{
			name:  "different year",
			td:    TransactionDate{Year: 2026, Month: 2, Day: 1},
			other: TransactionDate{Year: 2025, Month: 2, Day: 1},
			want:  false,
		},
		{
			name:  "different month",
			td:    TransactionDate{Year: 2026, Month: 2, Day: 1},
			other: TransactionDate{Year: 2026, Month: 3, Day: 1},
			want:  false,
		},
		{
			name:  "different day",
			td:    TransactionDate{Year: 2026, Month: 2, Day: 1},
			other: TransactionDate{Year: 2026, Month: 2, Day: 2},
			want:  false,
		},
		{
			name:  "both zero",
			td:    TransactionDate{},
			other: TransactionDate{},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.td.Equal(tt.other); got != tt.want {
				t.Errorf("TransactionDate.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		td      TransactionDate
		want    string
		wantErr bool
	}{
		{
			name: "normal date",
			td:   TransactionDate{Year: 2026, Month: 2, Day: 1},
			want: `"2026-02-01"`,
		},
		{
			name: "single digit month and day",
			td:   TransactionDate{Year: 2026, Month: 1, Day: 5},
			want: `"2026-01-05"`,
		},
		{
			name: "zero value",
			td:   TransactionDate{},
			want: "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.td.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionDate.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("TransactionDate.MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTransactionDate_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    TransactionDate
		wantErr bool
	}{
		{
			name: "normal date",
			json: `"2026-02-01"`,
			want: TransactionDate{Year: 2026, Month: 2, Day: 1},
		},
		{
			name: "single digit month and day",
			json: `"2026-01-05"`,
			want: TransactionDate{Year: 2026, Month: 1, Day: 5},
		},
		{
			name: "null value",
			json: "null",
			want: TransactionDate{},
		},
		{
			name: "empty string",
			json: `""`,
			want: TransactionDate{},
		},
		{
			name:    "invalid format",
			json:    `"2026/02/01"`,
			wantErr: true,
		},
		{
			name:    "invalid date",
			json:    `"2026-02-30"`,
			wantErr: true,
		},
		{
			name:    "not a string",
			json:    `{"year":2026,"month":2,"day":1}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TransactionDate
			err := json.Unmarshal([]byte(tt.json), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionDate.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("TransactionDate.UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionDate_JSONRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		td   TransactionDate
	}{
		{
			name: "normal date",
			td:   TransactionDate{Year: 2026, Month: 2, Day: 1},
		},
		{
			name: "leap year date",
			td:   TransactionDate{Year: 2024, Month: 2, Day: 29},
		},
		{
			name: "end of year",
			td:   TransactionDate{Year: 2025, Month: 12, Day: 31},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.td)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			// Unmarshal back
			var got TransactionDate
			err = json.Unmarshal(jsonData, &got)
			if err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			// Should equal original
			if !got.Equal(tt.td) {
				t.Errorf("Round trip failed: got %v, want %v", got, tt.td)
			}
		})
	}
}
