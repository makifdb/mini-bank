package models

import (
	"database/sql"
	"time"

	"github.com/vishal-bihani/go-tsid"
)

// Base struct defines the common fields for all entities.
type Base struct {
	ID         int64        `json:"-" pg:",pk"`                   // Internal unique identifier for the entity
	ExternalID string       `json:"id" pg:",unique,notnull"`      // External unique identifier for the entity
	CreatedAt  time.Time    `json:"created_at" pg:",notnull"`     // Timestamp when the entity was created
	UpdatedAt  time.Time    `json:"updated_at" pg:",notnull"`     // Timestamp when the entity was last updated
	DeletedAt  sql.NullTime `json:"deleted_at" pg:",soft_delete"` // Timestamp when the entity was deleted
}

// NewBase creates a new Base with the current time as CreatedAt and UpdatedAt,
// and generates a time-sorted unique identifier (TSID) for ExternalID.
func NewBase() Base {
	tsid := tsid.Fast() // Generate a fast TSID
	now := time.Now()   // Get the current time

	return Base{
		ID:         tsid.ToNumber(), // Convert TSID to a numeric value for the ID
		ExternalID: tsid.ToString(), // Convert TSID to a string for the ExternalID
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (b *Base) SetDeletedAt(t time.Time) {
	b.DeletedAt = sql.NullTime{Time: t, Valid: true}
}

func (b *Base) ClearDeletedAt() {
	b.DeletedAt = sql.NullTime{Valid: false}
}

func (b *Base) IsDeleted() bool {
	return b.DeletedAt.Valid
}
