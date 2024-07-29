package base

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common fields for all entities.
type Base struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`              // Unique identifier
	CreatedAt time.Time      `json:"created_at"`                                   // Timestamp when the entity was created
	UpdatedAt time.Time      `json:"updated_at"`                                   // Timestamp when the entity was last updated
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index" swaggertype:"string"` // Timestamp when the entity was deleted
}

// NewBase creates a new Base with the current time as CreatedAt and UpdatedAt,
// and generates a time-sorted unique identifier (TSID) for ID.
func NewBase() Base {
	now := time.Now()
	return Base{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
