package model

import (
	"time"
)

// Snippet represents a code snippet in the database.
type Snippet struct {
	ID        *string   `json:"id"`         // UUID string for the snippet's ID
	CreatedAt time.Time `json:"created_at"` // Timestamp of snippet creation
	Value     string    `json:"value"`      // The actual content of the snippet
	Name      string    `json:"name"`       // A name or title for the snippet
	Author    *string   `json:"author"`     // UUID string of the user who authored the snippet
}
