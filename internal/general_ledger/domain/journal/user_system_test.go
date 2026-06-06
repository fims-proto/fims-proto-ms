package journal

import (
	"testing"

	"github.com/google/uuid"
)

// TestSystemUserReviewBypassesSoD verifies SYSTEM_USER can review a journal it created.
// This tests the SoD bypass at the Review() method level.
func TestSystemUserReviewBypassesSoD(t *testing.T) {
	tests := []struct {
		name       string
		creator    string
		reviewer   string
		shouldFail bool
		errMsg     string
	}{
		{
			name:       "SYSTEM can review its own journal",
			creator:    SystemUser,
			reviewer:   SystemUser,
			shouldFail: false,
		},
		{
			name:       "Real user cannot review own journal",
			creator:    uuid.New().String(),
			reviewer:   "", // will be set to creator
			shouldFail: true,
			errMsg:     "journal-review-reviewerSameAsCreator",
		},
		{
			name:       "Different reviewer can review",
			creator:    uuid.New().String(),
			reviewer:   uuid.New().String(),
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviewer := tt.reviewer
			if reviewer == "" {
				reviewer = tt.creator // same as creator for real user test
			}

			// Create a minimal journal struct just for testing the Review logic
			j := &Journal{
				creator: tt.creator,
				auditor: "",
			}

			err := j.Review(reviewer)
			if tt.shouldFail && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldFail && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestSystemUserAuditBypassesSoD verifies SYSTEM_USER can audit a journal it created or reviewed.
func TestSystemUserAuditBypassesSoD(t *testing.T) {
	tests := []struct {
		name       string
		creator    string
		reviewer   string
		auditor    string
		shouldFail bool
	}{
		{
			name:       "SYSTEM can audit its own journal",
			creator:    SystemUser,
			reviewer:   SystemUser,
			auditor:    SystemUser,
			shouldFail: false,
		},
		{
			name:       "Real auditor cannot be creator",
			creator:    uuid.New().String(),
			reviewer:   uuid.New().String(),
			auditor:    "", // will be creator
			shouldFail: true,
		},
		{
			name:       "Different auditor can audit",
			creator:    uuid.New().String(),
			reviewer:   uuid.New().String(),
			auditor:    uuid.New().String(),
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auditor := tt.auditor
			if auditor == "" {
				auditor = tt.creator // same as creator for test
			}

			j := &Journal{
				creator:  tt.creator,
				reviewer: tt.reviewer,
			}

			err := j.Audit(auditor)
			if tt.shouldFail && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldFail && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestUpdateCheckBypassesForSystem verifies SYSTEM can update its own journals.
func TestUpdateCheckBypassesForSystem(t *testing.T) {
	tests := []struct {
		name       string
		creator    string
		updater    string
		shouldFail bool
	}{
		{
			name:       "SYSTEM can update own journal",
			creator:    SystemUser,
			updater:    SystemUser,
			shouldFail: false,
		},
		{
			name:       "Creator can update own journal",
			creator:    uuid.New().String(),
			updater:    "", // will be creator
			shouldFail: false,
		},
		{
			name:       "Non-creator cannot update",
			creator:    uuid.New().String(),
			updater:    uuid.New().String(),
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updater := tt.updater
			if updater == "" {
				updater = tt.creator
			}

			j := &Journal{
				creator: tt.creator,
			}

			err := j.checkUpdatePossible(updater)
			if tt.shouldFail && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.shouldFail && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
