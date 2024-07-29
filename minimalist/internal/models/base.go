package models

import (
	"time"
)

type Base struct {
	ID        int64      `json:"id"`         // Unique identifier
	CreatedAt time.Time  `json:"created_at"` // Timestamp when the entity was created
	UpdatedAt time.Time  `json:"updated_at"` // Timestamp when the entity was last updated
	DeletedAt *time.Time `json:"deleted_at"` // Timestamp when the entity was deleted, can be null
}

func NewBase() Base {
	now := time.Now()
	return Base{
		CreatedAt: now,
		UpdatedAt: now,
	}
}
