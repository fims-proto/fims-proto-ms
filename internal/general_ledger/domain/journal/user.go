package journal

import "github.com/google/uuid"

// SystemUser is the sentinel for system-automated actions (closing, auto-review, etc).
const SystemUser = "SYSTEM"

// emptyUser means "not yet assigned" (analogous to old uuid.Nil).
const emptyUser = ""

// SystemUserDBUUID is stored in the Postgres uuid column to represent SYSTEM.
// Distinct from uuid.Nil (= "not set"). Used only by the DB adapter.
var SystemUserDBUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

// IsSystemUser reports whether id represents the system actor.
func IsSystemUser(id string) bool { return id == SystemUser }

// isEmptyUser reports whether id is unset.
func isEmptyUser(id string) bool { return id == emptyUser }
