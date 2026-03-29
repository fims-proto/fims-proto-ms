package account

import (
	"testing"
)

func TestPadRawAccountNumber(t *testing.T) {
	tests := []struct {
		name      string
		hierarchy []int
		expected  string
		wantErr   bool
	}{
		{
			name:      "single level",
			hierarchy: []int{3104},
			expected:  "003104",
			wantErr:   false,
		},
		{
			name:      "two levels",
			hierarchy: []int{3104, 2},
			expected:  "003104000002",
			wantErr:   false,
		},
		{
			name:      "three levels",
			hierarchy: []int{3104, 2, 5},
			expected:  "003104000002000005",
			wantErr:   false,
		},
		{
			name:      "level with leading zeros in display",
			hierarchy: []int{1, 1},
			expected:  "000001000001",
			wantErr:   false,
		},
		{
			name:      "max levels (10)",
			hierarchy: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected:  "000001000002000003000004000005000006000007000008000009000010",
			wantErr:   false,
		},
		{
			name:      "empty hierarchy",
			hierarchy: []int{},
			wantErr:   true,
		},
		{
			name:      "exceeds max levels",
			hierarchy: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			wantErr:   true,
		},
		{
			name:      "negative segment",
			hierarchy: []int{1, -1},
			wantErr:   true,
		},
		{
			name:      "segment out of range",
			hierarchy: []int{1, 1_000_000},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PadRawAccountNumber(tt.hierarchy)
			if (err != nil) != tt.wantErr {
				t.Errorf("PadRawAccountNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("PadRawAccountNumber() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHierarchyFromRaw(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected []int
		wantErr  bool
	}{
		{
			name:     "single level",
			raw:      "003104",
			expected: []int{3104},
			wantErr:  false,
		},
		{
			name:     "two levels",
			raw:      "003104000002",
			expected: []int{3104, 2},
			wantErr:  false,
		},
		{
			name:     "three levels",
			raw:      "003104000002000005",
			expected: []int{3104, 2, 5},
			wantErr:  false,
		},
		{
			name:    "empty raw",
			raw:     "",
			wantErr: true,
		},
		{
			name:    "invalid length (not multiple of 6)",
			raw:     "0031040000021",
			wantErr: true,
		},
		{
			name:    "non-numeric segment",
			raw:     "00310400000a",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HierarchyFromRaw(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("HierarchyFromRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != len(tt.expected) {
					t.Errorf("HierarchyFromRaw() length = %d, want %d", len(got), len(tt.expected))
					return
				}
				for i, v := range got {
					if v != tt.expected[i] {
						t.Errorf("HierarchyFromRaw()[%d] = %d, want %d", i, v, tt.expected[i])
					}
				}
			}
		})
	}
}

func TestPadRawAccountNumberRoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		hierarchy []int
	}{
		{name: "single level", hierarchy: []int{3104}},
		{name: "two levels", hierarchy: []int{3104, 2}},
		{name: "three levels", hierarchy: []int{3104, 2, 5}},
		{name: "all ones", hierarchy: []int{1, 1, 1}},
		{name: "max values per level", hierarchy: []int{999999, 999999}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := PadRawAccountNumber(tt.hierarchy)
			if err != nil {
				t.Fatalf("PadRawAccountNumber() error = %v", err)
			}

			got, err := HierarchyFromRaw(raw)
			if err != nil {
				t.Fatalf("HierarchyFromRaw() error = %v", err)
			}

			if len(got) != len(tt.hierarchy) {
				t.Errorf("round-trip length mismatch: got %d, want %d", len(got), len(tt.hierarchy))
				return
			}
			for i, v := range got {
				if v != tt.hierarchy[i] {
					t.Errorf("round-trip[%d] = %d, want %d", i, v, tt.hierarchy[i])
				}
			}
		})
	}
}

func TestReadableFromRaw(t *testing.T) {
	tests := []struct {
		name        string
		raw         string
		codeLengths []int
		expected    string
		wantErr     bool
	}{
		{
			name:        "two levels with codeLengths [4,2]",
			raw:         "003104000002",
			codeLengths: []int{4, 2},
			expected:    "310402",
			wantErr:     false,
		},
		{
			name:        "three levels with codeLengths [4,2,2]",
			raw:         "003104000002000005",
			codeLengths: []int{4, 2, 2},
			expected:    "31040205",
			wantErr:     false,
		},
		{
			name:        "single level",
			raw:         "000001",
			codeLengths: []int{4},
			expected:    "0001",
			wantErr:     false,
		},
		{
			name:        "codeLengths shorter than hierarchy",
			raw:         "003104000002000005",
			codeLengths: []int{4, 2},
			wantErr:     true,
		},
		{
			name:        "invalid raw format",
			raw:         "0031040000021",
			codeLengths: []int{4, 2},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadableFromRaw(tt.raw, tt.codeLengths)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadableFromRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("ReadableFromRaw() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRawFromReadable(t *testing.T) {
	tests := []struct {
		name        string
		readable    string
		codeLengths []int
		expected    string
		wantErr     bool
	}{
		{
			name:        "two levels with codeLengths [4,2]",
			readable:    "310402",
			codeLengths: []int{4, 2},
			expected:    "003104000002",
			wantErr:     false,
		},
		{
			name:        "three levels with codeLengths [4,2,2]",
			readable:    "31040205",
			codeLengths: []int{4, 2, 2},
			expected:    "003104000002000005",
			wantErr:     false,
		},
		{
			name:        "single level",
			readable:    "0001",
			codeLengths: []int{4},
			expected:    "000001",
			wantErr:     false,
		},
		{
			name:        "level 1 only",
			readable:    "3104",
			codeLengths: []int{4, 2},
			expected:    "003104",
			wantErr:     false,
		},
		{
			name:        "level 1 and level 2",
			readable:    "310402",
			codeLengths: []int{4, 2, 2},
			expected:    "003104000002",
			wantErr:     false,
		},
		{
			name:        "code length not match",
			readable:    "3104021",
			codeLengths: []int{4, 2, 2},
			wantErr:     true,
		},
		{
			name:        "non-numeric segment",
			readable:    "310a02",
			codeLengths: []int{4, 2},
			wantErr:     true,
		},
		{
			name:        "length mismatch",
			readable:    "3104020",
			codeLengths: []int{4, 2},
			wantErr:     true,
		},
		{
			name:        "empty readable",
			readable:    "",
			codeLengths: []int{4},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RawFromReadable(tt.readable, tt.codeLengths)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawFromReadable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("RawFromReadable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRawFromReadableRoundTrip(t *testing.T) {
	tests := []struct {
		name        string
		readable    string
		codeLengths []int
	}{
		{name: "two levels [4,2]", readable: "310402", codeLengths: []int{4, 2}},
		{name: "three levels [4,2,2]", readable: "31040205", codeLengths: []int{4, 2, 2}},
		{name: "single level [4]", readable: "0001", codeLengths: []int{4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := RawFromReadable(tt.readable, tt.codeLengths)
			if err != nil {
				t.Fatalf("RawFromReadable() error = %v", err)
			}

			got, err := ReadableFromRaw(raw, tt.codeLengths)
			if err != nil {
				t.Fatalf("ReadableFromRaw() error = %v", err)
			}

			if got != tt.readable {
				t.Errorf("round-trip readable = %v, want %v", got, tt.readable)
			}
		})
	}
}

func TestLevelFromRaw(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected int
	}{
		{name: "one level", raw: "000001", expected: 1},
		{name: "two levels", raw: "000001000002", expected: 2},
		{name: "three levels", raw: "000001000002000003", expected: 3},
		{name: "ten levels", raw: "000001000002000003000004000005000006000007000008000009000010", expected: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LevelFromRaw(tt.raw)
			if got != tt.expected {
				t.Errorf("LevelFromRaw() = %d, want %d", got, tt.expected)
			}
		})
	}
}
