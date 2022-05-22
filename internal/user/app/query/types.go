package query

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID
	Traits    json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}
