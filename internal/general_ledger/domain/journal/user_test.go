package journal

import (
	"testing"

	"github.com/google/uuid"
)

func TestIsSystemUser(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{
			name:     "SYSTEM user",
			id:       SystemUser,
			expected: true,
		},
		{
			name:     "empty string",
			id:       "",
			expected: false,
		},
		{
			name:     "UUID string",
			id:       uuid.New().String(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSystemUser(tt.id)
			if got != tt.expected {
				t.Errorf("IsSystemUser(%q) = %v, want %v", tt.id, got, tt.expected)
			}
		})
	}
}

func TestIsEmptyUser(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{
			name:     "empty string",
			id:       "",
			expected: true,
		},
		{
			name:     "SYSTEM user",
			id:       SystemUser,
			expected: false,
		},
		{
			name:     "UUID string",
			id:       uuid.New().String(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isEmptyUser(tt.id)
			if got != tt.expected {
				t.Errorf("isEmptyUser(%q) = %v, want %v", tt.id, got, tt.expected)
			}
		})
	}
}

func TestSystemUserDBUUID(t *testing.T) {
	expected := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	if SystemUserDBUUID != expected {
		t.Errorf("SystemUserDBUUID = %v, want %v", SystemUserDBUUID, expected)
	}
}
