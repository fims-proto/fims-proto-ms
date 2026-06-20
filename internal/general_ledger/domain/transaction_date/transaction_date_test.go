package transaction_date

import "testing"

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
